# Auth Service – Production Readiness Questionnaire - Answers

Based on the code analysis, here are the answers to your production readiness questionnaire:

## 1. SYSTEM GOALS & THREAT MODEL

### 1.1 Business Context

* Apa skala target?

  * [x] side project
  * [ ] internal system
  * [ ] public SaaS
  * [ ] high-risk (finance / health / identity)

* Toleransi downtime?

  * [ ] < 1 menit
  * [ ] < 10 menit
  * [x] best effort

### 1.2 Threat Model (jawab eksplisit)

* Ancaman utama yang kamu lindungi:

  * [x] token theft
  * [x] replay attack
  * [x] brute force
  * [x] credential stuffing
  * [x] insider misuse

* Mana yang **tidak kamu lindungi secara sadar**?

  The service does not explicitly protect against:
  - DDoS attacks at the infrastructure level
  - Advanced persistent threats (APT)
  - Social engineering attacks
  - Database-level attacks if the database is compromised

---

## 2. TOKEN ARCHITECTURE

### 2.1 Access Token

* TTL: 30 minutes (configurable via JWT_EXPIRED environment variable)
* Algoritma: HS256 (HMAC with SHA-256)
* Claims wajib:

  * [x] sub
  * [x] iat
  * [x] exp
  * [ ] sid
  * [x] role
  * [ ] token_version

* Apakah access token bisa di-revoke secara soft?

  * Jika ya, bagaimana?
    Access tokens cannot be revoked individually since they are stateless JWTs. However, they have a short TTL (30 minutes) which limits exposure.

### 2.2 Refresh Token

* TTL hard limit: 7 days (configurable via JWT_REFRESH_EXPIRED environment variable)
* Rotation:

  * [x] enabled
  * [ ] disabled

* Reuse detection:

  * [x] implemented
  * [ ] belum

* Jika reuse terdeteksi, apa **respon sistem**?

  When reuse is detected, the system revokes the refresh token and returns an error "refresh token reuse detected".

---

## 3. REDIS DESIGN (STATE MANAGEMENT)

### 3.1 Key Strategy

Sebutkan semua key pattern:

* refresh:{jti}
* refresh:used:{jti} (for reuse detection)

TTL masing-masing key:
* refresh:{jti}: Same as refresh token TTL (7 days by default)
* refresh:used:{jti}: 2 minutes (to detect reuse)

### 3.2 Atomicity

* Apakah refresh flow atomic?

  * [ ] Lua script
  * [x] singleflight
  * [ ] tidak atomic

* Risiko race condition yang sudah kamu identifikasi?

  The refresh token rotation process includes steps to prevent race conditions by marking old tokens as used temporarily and removing them after rotation.

---

## 4. SESSION & DEVICE CONTROL

* Apakah ada konsep session / device?

  * [x] sid
  * [ ] device_id
  * [ ] tidak ada

* Apakah user bisa:

  * logout satu device? Yes, by revoking the refresh token
  * lihat daftar session? No, not implemented

* Data apa yang disimpan per session?

  * [ ] ip
  * [ ] ua
  * [x] created_at

---

## 5. FAILURE MODE & RESILIENCE

### 5.1 Redis Down Scenario

* Jika Redis **down total**:

  * Access token: Still work since they are stateless JWTs, but refresh tokens will fail
  * Refresh token: Will fail since they require Redis validation

* Timeout vs hard down dibedakan?

  The system checks Redis connectivity at startup and will fail to start if Redis is not available.

### 5.2 Clock & Time

* Clock skew tolerance: Not explicitly implemented, relies on JWT library defaults
* nbf digunakan? No, not currently implemented

---

## 6. RATE LIMIT & ABUSE CONTROL

Endpoint berikut, sebutkan limitnya:

* login: Rate limited via middleware (specific limits depend on implementation)
* refresh: Rate limited via middleware
* revoke: Rate limited via middleware

Apakah limit berbasis:

* [x] IP
* [ ] user
* [ ] device

---

## 7. OBSERVABILITY & AUDIT

### 7.1 Logging

Event apa saja yang dicatat?

* [x] login_success
* [x] login_failed
* [x] refresh_success
* [x] refresh_reuse
* [x] revoke_all

Struktur log:

* [x] JSON
* [ ] plain text

### 7.2 Monitoring

Metric apa yang kamu pantau?

* refresh per user
* reuse rate
* error rate

---

## 8. SECURITY BOUNDARY

* Apakah refresh token pernah:

  * lewat JS? No, should be stored in httpOnly cookies
  * lewat URL? No
  * disimpan di localStorage? No, should use httpOnly cookies

* Cookie policy:

  * HttpOnly? Not explicitly implemented in the code shown
  * SameSite? Not explicitly implemented in the code shown

---

## 9. MICROSERVICE READINESS

* Apakah service lain:

  * verify JWT sendiri? Yes, other services can verify JWTs using the shared secret
  * call auth service setiap request? No, they can verify tokens independently

* Apakah signing key:

  * [x] shared
  * [ ] rotated