# 🛰️ OpsWatch

**Real-time Linux server monitoring, right in your terminal.**

*Lightweight. Fast. Native. Zero fuss.*

![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/platform-Linux%20amd64-333333?logo=linux&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-brightgreen)
![Release](https://img.shields.io/badge/release-v0.1.5-blueviolet)

---

OpsWatch is a lightweight Linux CLI dashboard for monitoring server resources in real time — built in Go, shipped as a native Debian package, and designed to stay out of your way.

No web UI. No background daemon. No bloat. Just `ssh` in, run `opswatch`, and see what your server is doing.

---

## ✨ Features

| Feature | Description |
|---|---|
| 🖥️ **CPU Monitoring** | Real-time CPU usage |
| 🧠 **Memory Monitoring** | Live memory statistics |
| 💾 **Disk Usage** | Disk space monitoring |
| 📊 **History Chart** | Resource history inside the terminal |
| ⚡ **Lightweight** | Minimal system overhead |
| 🐧 **Native Linux** | No runtime dependencies |
| 📦 **Debian Package** | Easy installation via `apt` |

---

## 📥 Installation

### Ubuntu / Debian

Download the latest Debian package:

```bash
wget https://github.com/hrp7dev/opswatch/releases/latest/download/opswatch_amd64.deb
```

Install it:

```bash
sudo apt install ./opswatch_amd64.deb
```

Run it:

```bash
opswatch
```

That's it. 🎉

---

## 🚀 Usage

Start OpsWatch:

```bash
opswatch
```

The dashboard displays, live:

- 🖥️ CPU usage
- 🧠 Memory usage
- 💾 Disk usage
- 📊 Resource history

Exit anytime with:

```
Ctrl + C
```

---

## 📦 Releases

Every release contains:

- Linux amd64 Debian package (`.deb`)
- Automated GitHub Actions build

### 🛠️ Supported Systems

| OS | Status |
|---|---|
| Ubuntu 22.04+ | ✅ Supported |
| Ubuntu 24.04+ | ✅ Supported |
| Debian 12+ | ✅ Supported |

**Architecture:** `amd64`

---

## 🗺️ Roadmap

- [ ] Network monitoring
- [ ] Process monitoring
- [ ] Docker monitoring
- [ ] Alert system
- [ ] Web dashboard
- [ ] APT repository

---

## 👤 Author

**HamidReza Pouretemadi**
Fullstack Developer

- GitHub: [@hrp7dev](https://github.com/hrp7dev)
- Email: hrp7dev@gmail.com

---

## 📄 License

Released under the [MIT License](LICENSE).

---

*Made with ☕ and Go, for anyone who'd rather monitor a server than babysit one.*