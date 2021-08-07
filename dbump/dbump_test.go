package dbump

var testdataMigrations = []*Migration{
	&Migration{
		ID:       1,
		Apply:    `SELECT 1;`,
		Rollback: `SELECT 10;`,
	},
	&Migration{
		ID:       2,
		Apply:    `SELECT 2;`,
		Rollback: `SELECT 20;`,
	},
	&Migration{
		ID:       3,
		Apply:    `SELECT 3;`,
		Rollback: `SELECT 30;`,
	},
}
