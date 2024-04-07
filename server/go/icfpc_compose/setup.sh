ICFPC_BASE_PATH=/opt/icfpc

OVERWRITE=no
opts="$(getopt -o oph -l overwrite,help --name "$0" -- "$@")"
eval set -- "$opts"

while true
do
  case "$1" in
    -o|--overwrite)
      OVERWRITE=yes
      shift
      ;;
    -h|--help)
      echo " ICFPC setup script usage: $0 [-o|--overwrite]"
      echo "  The --overwrite option replaces any existing configuration with a fresh new install"
      exit 1
      ;;
    --)
      shift
      break
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done


create_env() {
  echo "** Creating icfpc environment file **"

  # Minimum mandatory values (when not just developing)
  COOKIE_SECRET=$(pwgen -1s 32)
  SECRET_KEY=$(pwgen -1s 32)
  PG_PASSWORD=$(pwgen -1s 32)
  DATABASE_URL="postgresql://postgres:${PG_PASSWORD}@icfpc_postgres/postgres?sslmode=disable"

  if [ -e "$ICFPC_BASE_PATH"/env ]; then
    # There's already an environment file

    if [ "x$OVERWRITE" = "xno" ]; then
      echo
      echo "Environment file already exists, reusing that one + and adding any missing (mandatory) values"

      POSTGRES_PASSWORD=
      POSTGRES_PASSWORD=$(. "$ICFPC_BASE_PATH"/env && echo "$POSTGRES_PASSWORD")
      if [ -z "$POSTGRES_PASSWORD" ]; then
        POSTGRES_PASSWORD=$PG_PASSWORD
        echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD" >> "$ICFPC_BASE_PATH"/env
        echo "POSTGRES_PASSWORD added to env file"
      fi

      DATABASE_URL=
      DATABASE_URL=$(. "$ICFPC_BASE_PATH"/env && echo "$DATABASE_URL")
      if [ -z "$DATABASE_URL" ]; then
        echo "DATABASE_URL=postgresql://postgres:${POSTGRES_PASSWORD}@icfpc_postgres/postgres?sslmode=disable" >> "$ICFPC_BASE_PATH"/env
        echo "DATABASE_URL added to env file"
      fi

      echo
      return
    fi

    # Move any existing environment file out of the way
    mv -f "$ICFPC_BASE_PATH"/env "$ICFPC_BASE_PATH"/env.old
  fi

  echo "Generating brand new environment file"

  cat <<EOF >"$ICFPC_BASE_PATH"/env
POSTGRES_PASSWORD=$PG_PASSWORD
DATABASE_URL=$DATABASE_URL
EOF
}

create_directories() {
  echo "** Creating $ICFPC_BASE_PATH directory structure for ICFPC **"

  if [ ! -e "$ICFPC_BASE_PATH" ]; then
    mkdir -p "$ICFPC_BASE_PATH"
    chown "$USER:" "$ICFPC_BASE_PATH"
  fi

  if [ -e "$ICFPC_BASE_PATH"/postgres-data ]; then
    # PostgreSQL database directory seems to exist already

    if [ "x$OVERWRITE" = "xyes" ]; then
      # We've been asked to overwrite the existing database, so delete the old one
      echo "Shutting down any running ICFPC instance"
      if [ -e "$ICFPC_BASE_PATH"/compose.yaml ]; then
        docker compose -f "$ICFPC_BASE_PATH"/compose.yaml down
      fi

      echo "Removing old ICFPC PG database directory"
      rm -rf "$ICFPC_BASE_PATH"/postgres-data
      mkdir "$ICFPC_BASE_PATH"/postgres-data
    fi
  else
    mkdir "$ICFPC_BASE_PATH"/postgres-data
  fi
}

if [ ! -d /opt/redash/ ]; then
    bash redash_setup/setup.sh
fi
create_directories
create_env
docker compose up --build