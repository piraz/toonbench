from toon_format import decode, encode


data = {
        "user": {
            "id": 123,
            "name": "Monstro do Lago",
            "active": True
        },
        "tags": ["python", "llm", "data"]
}

toon_output = encode(data)
print(toon_output)
data_from_toon = decode(toon_output)
print(data_from_toon)
