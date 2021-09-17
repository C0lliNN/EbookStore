echo "Starting localstack initialization..."

awslocal ses verify-domain-identity --domain ebook_store.com --endpoint-url http://localhost:4566

echo "Initialization has finished!"
