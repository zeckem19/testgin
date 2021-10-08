curl http://localhost:8000/register/K00001 \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data {"name": "Lim","age":40,"gender":"M"}