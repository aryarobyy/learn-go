FROM migrate/migrate

COPY ./db/migrations /db/migrations

ENTRYPOINT [ "migrate", "-path", "/migrations", "-database"]
CMD ["postgresql://postgres:123@auth.db:5432/authdb?sslmode=disable up"]
