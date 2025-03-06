sleep 5

echo "Opening ports..."

gh codespace ports visibility 5002:public -c $CODESPACE_NAME
gh codespace ports visibility 13008:public -c $CODESPACE_NAME
gh codespace ports visibility 15552:public -c $CODESPACE_NAME