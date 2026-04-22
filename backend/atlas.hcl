data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./internal/atlas/atlas.go",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "postgres://postgres:postgres@localhost:54322/postgres?search_path=public&sslmode=disable"
  migration {
    dir = "file://../supabase/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}