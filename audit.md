# AUDIT TEKNIS – GO AUTH SERVICE

*Baseline: Roadmap 4 Minggu – Mastering Go Auth Service (Advanced)*

Dokumen ini adalah **audit checklist formal** untuk menilai sejauh mana project auth service telah memenuhi **pedoman teknis** yang ditetapkan di roadmap 4 minggu. Dokumen ini **tidak memberikan solusi**, hanya mendefinisikan **apa saja yang harus diperiksa**.

---

## A. KONTEKS & TUJUAN SISTEM (Minggu 1)

### A1. Problem Statement

* Apakah auth service didefinisikan sebagai **read-heavy system**?
* Apakah target concurrency ditentukan secara eksplisit (angka)?
* Apakah SLA kasar didefinisikan (p95 / p99 latency)?

### A2. Scope & Batasan

* Apakah auth service berdiri sebagai **service terpisah**, bukan modul internal?
* Apakah scope auth dibatasi jelas pada:

  * login
  * validate token
  * refresh token
  * revoke token
* Apakah auth service **tidak menangani domain non-auth**?

---

## B. ARSITEKTUR & STRUKTUR KODE (Minggu 1)

### B1. Layering

* Apakah domain layer bebas dari dependency berikut:

  * HTTP framework
  * Redis client
  * JWT library
* Apakah service layer hanya mengorkestrasi flow?
* Apakah repository, cache, dan transport memiliki kontrak terpisah?

### B2. Dependency Direction

* Apakah dependency satu arah (tidak cyclic)?
* Apakah auth service dapat di-test tanpa menjalankan HTTP server?

---

## C. AUTH FLOW IMPLEMENTATION (Minggu 1)

### C1. Login Flow

* Apakah login flow sinkron dan deterministic?
* Apakah timeout diterapkan secara eksplisit?
* Apakah error DB tidak bocor ke client?

### C2. Validate Token Flow

* Apakah access token divalidasi secara stateless?
* Apakah validate token tidak melakukan call Redis pada kondisi normal?
* Apakah latency validate token rendah dan konsisten?

### C3. Refresh Token Flow

* Apakah refresh flow didefinisikan secara eksplisit?
* Apakah refresh token selalu diverifikasi ke Redis?
* Apakah refresh flow terpisah dari validate flow?

---

## D. CACHING & CONCURRENCY (Minggu 2)

### D1. Redis Usage

* Untuk apa saja Redis digunakan (state, cache, source of truth)?
* Apakah Redis boleh mati tanpa mematikan seluruh sistem?

### D2. Multi-layer Cache

* Apakah terdapat in-memory cache?
* Apakah Redis bukan satu-satunya lapisan cache?
* Apakah TTL antar layer cache berbeda dan disengaja?

### D3. Cache Stampede Protection

* Apakah `singleflight` digunakan?
* Digunakan pada flow apa saja?
* Apakah concurrency tinggi pernah disimulasikan?

---

## E. TTL, INVALIDATION & REVOCATION (Minggu 2)

* Apakah TTL access token short-lived?
* Apakah TTL konsisten dengan strategi cache?
* Apakah refresh token:

  * bisa di-revoke?
  * di-rotate?
* Apakah token lama benar-benar tidak usable setelah revoke?

---

## F. FAILURE HANDLING & TIMEOUT (Minggu 3)

### F1. Redis Failure Mode

* Apa perilaku sistem saat Redis:

  * down total?
  * timeout?
* Apakah access token tetap bisa divalidasi?

### F2. Timeout Strategy

* Apakah timeout diterapkan di:

  * HTTP layer
  * service layer
  * Redis
  * database
* Apakah timeout antar layer berbeda?

### F3. Graceful Shutdown

* Apakah SIGTERM ditangani?
* Apakah goroutine dihentikan dengan bersih?
* Apakah potensi goroutine leak dianalisis?

---

## G. OBSERVABILITY (Minggu 3)

### G1. Logging

* Apakah logging bersifat structured?
* Apakah setiap request memiliki request ID?
* Apakah log cukup untuk analisis kegagalan?

### G2. Metrics

* Apakah tersedia metrics untuk:

  * latency
  * error rate
  * cache hit/miss
* Apakah metrics mencerminkan dampak failure injection?

---

## H. FAILURE INJECTION (Minggu 3)

* Apakah Redis pernah dimatikan secara sengaja saat test?
* Apakah latency DB pernah disimulasikan?
* Apakah sistem tetap degrade secara terkontrol, bukan collapse?

---

## I. OPTIMIZATION & PEMBUKTIAN (Minggu 4)

### I1. Benchmark

* Apakah benchmark dibuat sebelum caching?
* Apakah benchmark dibuat sesudah caching?
* Apakah fokus benchmark pada latency dan allocation?

### I2. Profiling

* Apakah `pprof` digunakan untuk:

  * CPU
  * memory
  * allocation
* Apakah bottleneck diidentifikasi secara eksplisit?

### I3. Race Condition

* Apakah `go test -race` dijalankan?
* Apakah hasil race detector bersih?

---

## J. DOKUMENTASI & TRADE-OFF (Minggu 4)

* Apakah keputusan desain utama didokumentasikan?
* Apakah trade-off teknis dicatat secara eksplisit?
* Apakah ada bagian yang mengakui keterbatasan desain saat ini?

---

## PENUTUP

Audit ini bertujuan memastikan bahwa auth service:

* dapat dipertanggungjawabkan secara teknis
* tahan terhadap failure dan concurrency
* sesuai dengan tujuan pembelajaran roadmap

Dokumen ini digunakan sebagai **dasar review teknis**, bukan sebagai checklist fitur.
