echo "Starting localstack initialization..."

awslocal ses verify-domain-identity --domain ebook_store.com --endpoint-url http://localhost:4566
awslocal s3api create-bucket --bucket ebook-store --endpoint-url http://localhost:4566 --create-bucket-configuration LocationConstraint=us-east-2

echo "Initialization has finished!"
