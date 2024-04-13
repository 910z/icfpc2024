OVERWRITE_ENV=no
OVERWRITE_DB=no
ENV_NAME=dev
opts="$(getopt -o oph -l overwrite,help,prod,overwrite-db --name "$0" -- "$@")"
eval set -- "$opts"

while true
do
  case "$1" in
    -o|--overwrite)
      OVERWRITE_ENV=yes
      OVERWRITE_DB=yes
      shift
      ;;
    -o|--overwrite-db)
      OVERWRITE_DB=yes
      shift
      ;;
    --prod)
      ENV_NAME=prod
      shift
      ;;
    -h|--help)
      echo " ICFPC setup script usage: $0 [-o|--overwrite] [--prod]"
      echo "  The --overwrite option replaces any existing configuration with a fresh new install"
      echo "  The --overwrite-db option purges only db state, leaving environment file intact"
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

  PG_PASSWORD=$(pwgen -1s 32)
  
  if [ -e $SECRET_ENV_PATH ]; then
    # There's already an environment file

    if [ "x$OVERWRITE_ENV" = "xno" ]; then
      echo
      echo "Environment file already exists, reusing that one + and adding any missing (mandatory) values"

      POSTGRES_PASSWORD=
      POSTGRES_PASSWORD=$(. $SECRET_ENV_PATH && echo "$POSTGRES_PASSWORD")
      if [ -z "$POSTGRES_PASSWORD" ]; then
        POSTGRES_PASSWORD=$PG_PASSWORD
        echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD" >> $SECRET_ENV_PATH
        echo "POSTGRES_PASSWORD added to env file"
      fi
      return
    fi

    # Move any existing environment file out of the way
    mv -f $SECRET_ENV_PATH $SECRET_ENV_PATH.old
  fi

  echo "Generating brand new environment file"

  cat <<EOF >$SECRET_ENV_PATH
POSTGRES_PASSWORD=$PG_PASSWORD
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

    if [ "x$OVERWRITE_DB" = "xyes" ]; then
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
ENV_PATH=./$ENV_NAME.env
echo Reading env from $ENV_PATH
export $(cat $ENV_PATH | xargs) # переменные отсюда будут нужны и в функциях ниже, и в docker compose
create_directories
SECRET_ENV_PATH="$ICFPC_BASE_PATH"/secret.env
echo Creating secrets in $SECRET_ENV_PATH
create_env

detach=()
if [ "$ENV_NAME" = "prod" ]; then
   detach+=(--detach)
fi

docker compose --env-file $SECRET_ENV_PATH --project-name $ENV_NAME up --build "${detach[@]}"