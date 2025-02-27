cd /workspaces

git clone https://github.com/stakwork/sphinx-tribes-frontend 

cd sphinx-tribes

echo "DATABASE_URL=postgres://test_user:test_password@localhost:5432/test_db" > .env
echo "LN_JWT_KEY=notasecretstring" >> .env