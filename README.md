# backup-cli

A simple command-line tool written in Go that allows you to create ZIP archives of files and encrypt them using AES-256 (AES-GCM). Also supports decryption.

## 🔐 Features

- 🔒 AES-256 encryption (GCM mode)
- 📦 Archive multiple files into a single `.zip`
- 🔓 Decrypt and extract encrypted backup files
- 🧪 Minimal, dependency-free, written in pure Go

## 📦 Installation

Make sure you have Go installed:  
👉 https://golang.org/dl/

Then build the program:

```bash
go build -o backup-cli main.go
```

## 🚀 Usage

### Encrypt files

```bash
go run main.go -key "your-32-byte-secret-key-123456789012" file1.txt file2.txt
```

- This will create a `backup.zip.aes` encrypted archive.
- Replace files with any files you want to back up.

### Decrypt

```bash
go run main.go -key "your-32-byte-secret-key-123456789012" -decrypt -out restored.zip backup.zip.aes
```

- This will decrypt the encrypted backup and save it as `restored.zip`.

## 🧪 Testing (Optional)

To generate some test `.txt` files, you can write a helper Go script, or just manually create files like:

```bash
echo "Hello from file 1" > file1.txt
echo "This is file 2" > file2.txt
```

## ⚠️ Notes

- The encryption key must be **exactly 32 bytes** (256 bits) long.
- Encrypted output uses AES-GCM with a randomly generated nonce.

## 📄 License

MIT — free to use, modify, and share.

## 👤 Author

Made by [xyberis](https://github.com/TheXyberis)