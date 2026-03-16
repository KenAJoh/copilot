---
name: security-review
description: Bruk før commit, push eller pull request for å sjekke at koden er trygg å merge
---

# Security Review Skill

This skill provides pre-commit and pre-PR security checks for Nav applications. Covers secret scanning, vulnerability scanning, and Nav-specific requirements.

For architecture questions, threat modeling, or compliance decisions, use `@security-champion` instead.

## Automated Scans

Run with `run_in_terminal`:

```bash
# Scan repo for known vulnerabilities and secrets
trivy repo .

# Scan Docker image for HIGH/CRITICAL CVEs
trivy image <image-name> --severity HIGH,CRITICAL

# Scan GitHub Actions workflows for insecure patterns
zizmor .github/workflows/

# Quick search for secrets in git history
git log -p --all -S 'password' -- '*.kt' '*.ts' | head -100
git log -p --all -S 'secret' -- '*.kt' '*.ts' | head -100
```

## Parameterized SQL (Never Concatenate)

```kotlin
// ✅ Correct – parameterized query
fun findBruker(fnr: String): Bruker? =
    jdbcTemplate.queryForObject(
        "SELECT * FROM bruker WHERE fnr = ?",
        brukerRowMapper,
        fnr
    )

// ❌ Wrong – SQL injection risk
fun findBrukerUnsafe(fnr: String): Bruker? =
    jdbcTemplate.queryForObject(
        "SELECT * FROM bruker WHERE fnr = '$fnr'",
        brukerRowMapper
    )
```

## No PII in Logs

```kotlin
// ✅ Correct – log correlation ID, not PII
log.info("Behandler sak for bruker", kv("sakId", sak.id), kv("tema", sak.tema))

// ❌ Wrong – never log FNR, name, or other PII
log.info("Behandler sak for bruker ${bruker.fnr}")  // GDPR violation
log.info("Navn: ${bruker.navn}")                      // GDPR violation
```

## Secrets from Environment, Never Hardcoded

```kotlin
// ✅ Correct – read from environment (Nais injects via Secret)
val dbPassword = System.getenv("DB_PASSWORD")
    ?: throw IllegalStateException("DB_PASSWORD mangler")

// ❌ Wrong – hardcoded secret
val dbPassword = "supersecret123"
```

## Network Policy (Nais)

Only expose what must be exposed:

```yaml
spec:
  accessPolicy:
    inbound:
      rules:
        - application: frontend-app      # only explicitly named callers
    outbound:
      rules:
        - application: pdl-api
          namespace: pdl
          cluster: prod-gcp
      external:
        - host: api.external-service.no  # only if strictly necessary
```

## OWASP Top 10 Checks

### A01: Broken Access Control

```kotlin
// ✅ Korrekt — sjekk at bruker har tilgang til ressursen
@GetMapping("/api/vedtak/{id}")
fun getVedtak(@PathVariable id: UUID): ResponseEntity<VedtakDTO> {
    val bruker = hentInnloggetBruker()
    val vedtak = vedtakService.findById(id)
    if (vedtak.brukerId != bruker.id) {
        return ResponseEntity.status(HttpStatus.FORBIDDEN).build()
    }
    return ResponseEntity.ok(vedtak.toDTO())
}

// ❌ Feil — ingen tilgangskontroll (IDOR)
@GetMapping("/api/vedtak/{id}")
fun getVedtak(@PathVariable id: UUID) = vedtakService.findById(id)
```

### A03: Injection

```kotlin
// ✅ Korrekt — parameterisert spørring
jdbcTemplate.query("SELECT * FROM bruker WHERE fnr = ?", mapper, fnr)

// ❌ Feil — strengkonkatenering
jdbcTemplate.query("SELECT * FROM bruker WHERE fnr = '$fnr'", mapper)
```

### A05: Security Misconfiguration

```kotlin
// ✅ Korrekt — CORS kun for kjente domener
@Bean
fun corsFilter() = CorsFilter(CorsConfiguration().apply {
    allowedOrigins = listOf("https://my-app.intern.nav.no")
    allowedMethods = listOf("GET", "POST")
    allowedHeaders = listOf("Authorization", "Content-Type")
})

// ❌ Feil — åpen CORS
allowedOrigins = listOf("*")
```

### A07: Cross-Site Scripting (XSS)

```tsx
// ✅ Korrekt — React escaper automatisk
<BodyShort>{bruker.navn}</BodyShort>

// ❌ Feil — rå HTML-injeksjon
<div dangerouslySetInnerHTML={{ __html: userInput }} />
```

### A08: Insecure Deserialization

```kotlin
// ✅ Korrekt — valider input etter deserialisering
@PostMapping("/api/vedtak")
fun create(@RequestBody @Valid request: CreateVedtakRequest): ResponseEntity<VedtakDTO>

// ✅ Begrens Jackson til kjente typer
objectMapper.apply {
    activateDefaultTyping(
        polymorphicTypeValidator,
        ObjectMapper.DefaultTyping.NON_FINAL
    )
}
```

### A09: Logging & Monitoring

```kotlin
// ✅ Korrekt — strukturert logging med korrelerings-ID, ingen PII
log.info("Vedtak opprettet", kv("vedtakId", vedtak.id), kv("sakId", sak.id))

// ❌ Feil — PII i logger
log.info("Vedtak for bruker ${bruker.fnr} opprettet")
```

## File Upload Security

```kotlin
// ✅ Korrekt — valider filtype, størrelse, og magic bytes
fun validateUpload(file: MultipartFile) {
    require(file.size <= 10 * 1024 * 1024) { "Filen er for stor (maks 10 MB)" }
    require(file.contentType in ALLOWED_TYPES) { "Ugyldig filtype" }

    val bytes = file.bytes.take(8).toByteArray()
    require(verifyMagicBytes(bytes, file.contentType!!)) { "Filinnhold matcher ikke type" }
}

private val ALLOWED_TYPES = setOf("application/pdf", "image/png", "image/jpeg")
```

## Dependency Management

```kotlin
// build.gradle.kts — pin versjoner, bruk BOM
dependencyManagement {
    imports {
        mavenBom("org.springframework.boot:spring-boot-dependencies:3.4.1")
    }
}

// Sjekk sårbare avhengigheter
// ./gradlew dependencyCheckAnalyze
// trivy repo .
```

## Expanded Checklist

- [ ] SQL-spørringer er parameteriserte (ingen strengkonkatenering)
- [ ] Ingen PII i logger (fnr, navn, adresse)
- [ ] Hemmeligheter kun fra environment/secrets
- [ ] Nais accessPolicy er eksplisitt (ingen åpen inbound)
- [ ] CORS er begrenset til kjente domener
- [ ] Input er validert og sanitert
- [ ] Tilgangskontroll sjekker eierskap (ikke bare auth)
- [ ] Filopplasting validerer type, størrelse, og innhold
- [ ] Avhengigheter er oppdaterte og sårbarhetsskannet
- [ ] Ingen `dangerouslySetInnerHTML` uten sanitering

## Dependency Management

```bash
# Kotlin – check for outdated/vulnerable dependencies
./gradlew dependencyUpdates
./gradlew dependencyCheckAnalyze   # OWASP check

# Node/TypeScript
npm audit
npm audit fix
```

## Security Checklist

- [ ] No secrets, tokens, or API keys hardcoded in source
- [ ] No PII (FNR, name, address) in log statements
- [ ] All SQL queries use parameterized statements
- [ ] Nais `accessPolicy` limits inbound/outbound to only what is needed
- [ ] Token validation on all protected endpoints (see `@security-champion`)
- [ ] `trivy repo .` passes without HIGH/CRITICAL findings
- [ ] `zizmor` passes on all GitHub Actions workflows
- [ ] Git history clean of committed secrets (`git log` scan above)
- [ ] HTTPS enforced – no plain HTTP calls to external services
- [ ] Dependencies up to date (`dependencyUpdates` / `npm audit`)

## Related

| Resource | Use For |
|----------|---------|
| `@security-champion` | Threat modeling, compliance questions, Nav security architecture |
| `@auth-agent` | JWT validation, TokenX, ID-porten, Maskinporten |
| `@nais-agent` | Nais manifest, accessPolicy, secrets setup |
| [sikkerhet.nav.no](https://sikkerhet.nav.no) | Nav Golden Path, authoritative security guidance |
