## How are messages authenticated on client

1. When user sends a message to user2, he also sends the signature of pre shared message ("AUTHENTICITY VERIFICATION")
1. When user2 receives message, he verifies the signature using the public key of user (they have it in their local contact DB)
1. If signature verification fails, user2 discards the message

## How are users authenticated on server when user fetches unread messages

1. User sends signature of pre-shared message ("AUTHENTICITY VERIFICATION")
1. Server tries to verify the signature with public key mapped to the temporarily stored message
1. On success, server returns messages and deletes them from DB
1. Otherwise server doesn't do anything