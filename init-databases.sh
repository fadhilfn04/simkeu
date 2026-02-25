#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE simkeu_auth;
  CREATE DATABASE simkeu_blockchain;
  CREATE DATABASE simkeu_debitur;
  CREATE DATABASE simkeu_log;
  CREATE DATABASE simkeu_master;
  CREATE DATABASE simkeu_payment;
  CREATE DATABASE simkeu_piutang;
  CREATE DATABASE simkeu_realisasi;
  CREATE DATABASE simkeu_tagihan;
EOSQL
