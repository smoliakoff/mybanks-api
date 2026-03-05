env "local" {
  src = "ent://ent/schema"

  dev = "docker://postgres/16/dev"

  migration {
    dir = "file://migrations"
  }

  revisions_schema = "public"

  url = getenv("DATABASE_URL")

  schemas = ["public"]
}