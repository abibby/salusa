package migrate

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/database/schema"
	"github.com/abibby/salusa/set"
)

type DBMigration struct {
	model.BaseModel
	Name  string `db:"name,primary"`
	Run   bool   `db:"run"`
	table string
}

func (m *DBMigration) Table() string {
	return m.table
}

type Migrations struct {
	table      string
	migrations []*Migration
}

func New() *Migrations {
	return &Migrations{
		table:      "migrations",
		migrations: []*Migration{},
	}
}

func (m *Migrations) Add(migration *Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrations) isTableCreated(table string) bool {
	for _, m := range m.migrations {
		blueprinter, ok := m.Up.(schema.Blueprinter)
		if !ok {
			continue
		}
		if blueprinter.Type() == schema.BlueprintTypeCreate {
			if blueprinter.GetBlueprint().TableName() == table {
				return true
			}
		}
	}

	return false
}

func (m *Migrations) GenerateMigration(migrationName, packageName string, model model.Model) (string, error) {
	if !m.isTableCreated(database.GetTable(model)) {
		up, err := CreateFromModel(model)
		if err != nil {
			return "", err
		}
		return SrcFile(migrationName, packageName, up, drop(model))
	}

	up, down, err := m.update(model)
	if err != nil {
		return "", err
	}
	return SrcFile(migrationName, packageName, up, down)
}

func (m *Migrations) Blueprint(tableName string) *schema.Blueprint {
	result := &schema.Blueprint{}

	for _, migration := range m.migrations {
		blueprinter, ok := migration.Up.(schema.Blueprinter)
		if !ok {
			continue
		}
		blueprint := blueprinter.GetBlueprint()
		if blueprint.TableName() != tableName {
			continue
		}

		if blueprinter.Type() == schema.BlueprintTypeCreate {
			result = blueprint
		} else {
			result.Merge(blueprint)
		}
	}
	return result
}

func (m *Migrations) Up(ctx context.Context, db database.DB) error {
	sql, bindings, err := schema.Create(m.table, func(b *schema.Blueprint) {
		b.String("name")
		b.Bool("run")
	}).IfNotExists().ToSQL(dialects.New())
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, sql, bindings...)
	if err != nil {
		return err
	}

	migrations, err := builder.New[*DBMigration]().
		Select("*").
		From(m.table).
		OrderBy("name").
		WithContext(ctx).
		Get(db)
	if err != nil {
		return err
	}

	runMigrations := set.New[string]()
	for _, migration := range migrations {
		if migration.Run {
			runMigrations.Add(migration.Name)
		}
	}

	for _, migration := range m.migrations {
		if runMigrations.Has(migration.Name) {
			continue
		}

		m := &DBMigration{
			Name:  migration.Name,
			Run:   false,
			table: m.table,
		}
		err = model.SaveContext(ctx, db, m)
		if err != nil {
			return err
		}
		err := migration.Up.Run(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to prepare migration %s: %w", migration.Name, err)
		}

		m.Run = true
		err = model.SaveContext(ctx, db, m)
		if err != nil {
			return err
		}
	}
	return nil
}
