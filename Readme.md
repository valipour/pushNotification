# 📦 WebSocket Name Notifier – GoLang

یک سرور سبک با Go که از WebSocket استفاده می‌کند تا وقتی نام کاربر با `POST` تغییر می‌کند، فوراً به تمام کلاینت‌های متصل نوتیف ارسال کند.

بدون حلقه و بدون مصرف اضافی CPU — فقط push هنگام نیاز.

---

## 🚀 اجرای پروژه

### 1. نصب پکیج مورد نیاز
```bash
go get github.com/gorilla/websocket
```

### 2. اجرای برنامه
```bash
go run main.go
```

سرور روی آدرس زیر اجرا می‌شود:
```
http://localhost:8080
```

---

## 🧪 تست با Postman

### اتصال WebSocket:
1. وارد Postman شوید.
2. New → WebSocket Request
3. آدرس:
   ```
   ws://localhost:8080/ws
   ```
4. Connect را بزنید.

---

### ارسال نام جدید (HTTP POST)

```http
POST http://localhost:8080/set-name?name=علی
```

در پاسخ WebSocket:
```
نام جدید تنظیم شد: علی
```

---

## 🧑‍💻 کلاینت WebSocket برای تست (HTML/JS)

```js
const ws = new WebSocket("ws://localhost:8080/ws");
ws.onmessage = (e) => console.log("📢 نوتیف:", e.data);
```

---

## 📂 مسیرهای سرور

| مسیر         | روش   | عملکرد                                  |
|--------------|--------|-------------------------------------------|
| `/ws`        | GET (WebSocket) | اتصال WebSocket                       |
| `/set-name`  | POST   | تغییر نام و ارسال نوتیف به کلاینت‌ها     |

---

## ✅ ویژگی‌ها

- فقط ارسال پیام هنگام تغییر  نام
- بدون استفاده از حلقه یا زمان‌سنج
- مصرف بسیار کم منابع سیستم
- WebSocket برای ارتباط زنده
- مدیریت کلاینت‌ها با sync.Mutex

---

