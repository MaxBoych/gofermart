package sql_queries

import "fmt"

const (
	UserTableName = "user_account"

	UserIdColumnName   = "user_id"
	LoginColumnName    = "login"
	PasswordColumnName = "password"

	SecretKeyConstantColumnName = "constant_id"
	SecretKeyTableName          = "secret_key"
	SecretKeyValueColumnName    = "key_value"

	TokenTableName       = "token"
	TokenIdColumnName    = "token_id"
	TokenValueColumnName = "value"

	OrderTableName         = "user_order"
	OrderIdColumnName      = "order_id"
	OrderNumberColumnName  = "number"
	OrderAccrualColumnName = "accrual"
	OrderStatusColumnName  = "status"

	BalanceTableName           = "user_balance"
	BalanceIdColumnName        = "balance_id"
	BalanceCurrentColumnName   = "current"
	BalanceWithdrawnColumnName = "withdrawn"

	WithdrawTableName       = "user_withdraw"
	WithdrawIdColumnName    = "withdraw_id"
	WithdrawOrderColumnName = "order"
	WithdrawSumColumnName   = "sum"

	CreatedAtColumnName = "created_at"
	UpdatedAtColumnName = "updated_at"
	DeletedAtColumnName = "deleted_at"
)

var (
	SelectUser = []string{
		UserIdColumnName,
		LoginColumnName,
		PasswordColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
		DeletedAtColumnName,
	}

	InsertUser = []string{
		LoginColumnName,
		PasswordColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	InsertSecretKey = []string{
		SecretKeyConstantColumnName,
		SecretKeyValueColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	SelectSecretKey = []string{
		SecretKeyValueColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	InsertToken = []string{
		UserIdColumnName,
		TokenValueColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	SelectToken = []string{
		TokenIdColumnName,
		UserIdColumnName,
		TokenValueColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	InsertOrder = []string{
		OrderNumberColumnName,
		UserIdColumnName,
		OrderStatusColumnName,
		OrderAccrualColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	SelectOrder = []string{
		OrderIdColumnName,
		OrderNumberColumnName,
		UserIdColumnName,
		OrderStatusColumnName,
		OrderAccrualColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	InsertBalance = []string{
		UserIdColumnName,
		BalanceCurrentColumnName,
		BalanceWithdrawnColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	SelectBalance = []string{
		BalanceIdColumnName,
		UserIdColumnName,
		BalanceCurrentColumnName,
		BalanceWithdrawnColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	}

	InsertWithdraw = []string{
		WithdrawOrderColumnName,
		WithdrawSumColumnName,
		UserIdColumnName,
		CreatedAtColumnName,
	}

	SelectWithdraw = []string{
		WithdrawIdColumnName,
		WithdrawOrderColumnName,
		WithdrawSumColumnName,
		UserIdColumnName,
		CreatedAtColumnName,
	}

	SecretKeyTableSQL = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s" (
		    "%s" INT NOT NULL UNIQUE,
    		"%s" TEXT NOT NULL,
    		"%s" TIMESTAMP WITH TIME ZONE NOT NULL,
    		"%s" TIMESTAMP WITH TIME ZONE NOT NULL
		);`,
		SecretKeyTableName,
		SecretKeyConstantColumnName,
		SecretKeyValueColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	)

	JwtTableSQL = fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s" (
    		"%s" BIGSERIAL PRIMARY KEY,
    		"%s" BIGINT NOT NULL UNIQUE,
    		"%s" TEXT NOT NULL,
    		"%s" TIMESTAMP WITH TIME ZONE NOT NULL,
    		"%s" TIMESTAMP WITH TIME ZONE NOT NULL
		);`,
		TokenTableName,
		TokenIdColumnName,
		UserIdColumnName,
		TokenValueColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	)

	UserTableSQL = fmt.Sprintf(`
    	CREATE TABLE IF NOT EXISTS "%s" (
        	"%s" BIGSERIAL PRIMARY KEY,
        	"%s" TEXT NOT NULL UNIQUE,
        	"%s" TEXT NOT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE NOT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE NOT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE DEFAULT NULL
    	);`,
		UserTableName,
		UserIdColumnName,
		LoginColumnName,
		PasswordColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
		DeletedAtColumnName,
	)

	OrderTableSQL = fmt.Sprintf(`
    	CREATE TABLE IF NOT EXISTS "%s" (
        	"%s" BIGSERIAL PRIMARY KEY,
        	"%s" TEXT NOT NULL UNIQUE,
        	"%s" BIGINT NOT NULL,
        	"%s" TEXT NOT NULL,
        	"%s" DOUBLE PRECISION DEFAULT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE NOT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE NOT NULL
    	);`,
		OrderTableName,
		OrderIdColumnName,
		OrderNumberColumnName,
		UserIdColumnName,
		OrderStatusColumnName,
		OrderAccrualColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	)

	BalanceTableSQL = fmt.Sprintf(`
    	CREATE TABLE IF NOT EXISTS "%s" (
        	"%s" BIGSERIAL PRIMARY KEY,
        	"%s" BIGINT NOT NULL,
        	"%s" DOUBLE PRECISION NOT NULL,
        	"%s" DOUBLE PRECISION NOT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE NOT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE NOT NULL
    	);`,
		BalanceTableName,
		BalanceIdColumnName,
		UserIdColumnName,
		BalanceCurrentColumnName,
		BalanceWithdrawnColumnName,
		CreatedAtColumnName,
		UpdatedAtColumnName,
	)

	WithdrawTableSQL = fmt.Sprintf(`
    	CREATE TABLE IF NOT EXISTS "%s" (
        	"%s" BIGSERIAL PRIMARY KEY,
        	"%s" TEXT NOT NULL,
        	"%s" DOUBLE PRECISION NOT NULL,
        	"%s" BIGINT NOT NULL,
        	"%s" TIMESTAMP WITH TIME ZONE NOT NULL
    	);`,
		WithdrawTableName,
		WithdrawIdColumnName,
		WithdrawOrderColumnName,
		WithdrawSumColumnName,
		UserIdColumnName,
		CreatedAtColumnName,
	)
)
