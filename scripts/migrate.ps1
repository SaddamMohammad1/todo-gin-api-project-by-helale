# ------------------------------------------
# Load all environment variables from .env
# ------------------------------------------
# Read each line of .env file
Get-Content .env | ForEach-Object {
    # Match key=value format and ignore commented lines (#)
    if ($_ -match '^([^#][^=]+)=(.+)$') {
        # Set the environment variable inside PowerShell session
        # Example: DATABASE_URL=postgres://... → $env:DATABASE_URL
        Set-Item -Path "env:$($matches[1])" -Value $matches[2]
    }
}

# ------------------------------------------
# Read command-line arguments
# ------------------------------------------
# $args[0] = migration action (up, down, create, force)
# $args[1] = second argument (name or count)
$command = $args[0]
$name = $args[1]

# ------------------------------------------
# Handle migration commands using switch-case
# ------------------------------------------
switch ($command) {

    # ---------------------------
    # Run all pending migrations
    # Command: ./migrate.ps1 up
    # ---------------------------
    "up" {
        migrate -path migrations -database $env:DATABASE_URL up
    }

    # -----------------------------------------
    # Rollback migration(s)
    # Command: ./migrate.ps1 down
    # Command: ./migrate.ps1 down 3  (rollback 3)
    # -----------------------------------------
    "down" {

        # If user passed rollback count → use it
        # Otherwise rollback only 1 migration
        $count = if ($name) { $name } else { "1" }

        # Ask confirmation before rollback
        Write-Host "Rolling back $count migration(s). Continue? [y/N]"
        $confirm = Read-Host

        # Run rollback only when user types 'y'
        if ($confirm -eq 'y') {
            migrate -path migrations -database $env:DATABASE_URL down $count
        }
    }

    # -----------------------------------------------------------------
    # Create new migration files (up & down)
    # Command: ./migrate.ps1 create create_users_table
    # Output: migrations/000001_create_users_table.up.sql
    #         migrations/000001_create_users_table.down.sql
    # -----------------------------------------------------------------
    "create" {
        migrate create -ext sql -dir migrations -seq $name
    }

    # ------------------------------------------------------------
    # Force set migration version manually
    # Command: ./migrate.ps1 force 2
    # Used when migration version gets stuck or corrupted
    # ------------------------------------------------------------
    "force" {
        migrate -path migrations -database $env:DATABASE_URL force $name
    }
}
