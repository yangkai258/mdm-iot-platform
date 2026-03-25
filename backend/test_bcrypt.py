import bcrypt

# Generate a new hash
password = b'admin123'
new_hash = bcrypt.hashpw(password, bcrypt.gensalt())
print('New hash:', new_hash.decode())

# Verify
result = bcrypt.checkpw(password, new_hash)
print('Verify result:', result)
