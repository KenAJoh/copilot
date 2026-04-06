---
title: "Nyheter og trender — April 2026"
date: 2026-04-06
draft: true
category: copilot
excerpt: "Copilot SDK i public preview, legacy metrics API nedlagt, organisasjonsstyrt runner, personvernpolicy trer i kraft 24. april."
tags:
  - copilot-sdk
  - coding-agents
  - enterprise-controls
  - privacy
  - metrics
---

<!-- AI-REDAKSJONELT: Denne artikkelen er en oppsummering av de viktigste endringene og trendene — ikke en komplett liste. Prioriter det som er mest relevant for Nav-utviklere. Mindre oppdateringer samles i «Flere oppdateringer»-seksjonen. Individuelle nyheter dekkes av egne excerpt-filer i samme mappe. -->

April 2026 starter med infrastruktur. GitHub åpner Copilot-motoren som SDK, legger ned det gamle metrics-API-et, og gir organisasjoner bedre kontroll over hvordan coding agent kjører. Senere i måneden trer den kontroversielle personvernpolicyen for treningsdata i kraft.

---

## 1. Copilot SDK i public preview

GitHub Copilot SDK er nå tilgjengelig i public preview — det samme agentmotoren som driver Copilot cloud agent og Copilot CLI, pakket som bibliotek. SDK-et gir verktøyinvokning, streaming, filoperasjoner og multi-turn-sesjoner rett ut av boksen, uten at du trenger å bygge egen AI-orkestrering.

Tilgjengelig i fem språk: Node.js/TypeScript, Python, Go, .NET og Java (nytt). Nøkkelfunksjoner inkluderer custom tools med handlers, finkornet system-prompt-tilpasning (`replace`, `append`, `prepend`, `transform`), OpenTelemetry-integrasjon for distribuert tracing, et tillatelsesrammeverk for sensitive operasjoner, og Bring Your Own Key (BYOK) for OpenAI, Azure AI Foundry eller Anthropic.

SDK-et er tilgjengelig for alle — også brukere uten Copilot-abonnement via BYOK. Hver prompt teller mot premium request-kvoten for Copilot-abonnenter.

**Kilde:** [Copilot SDK in public preview](https://github.blog/changelog/2026-04-02-copilot-sdk-in-public-preview/) (GitHub Changelog, 2. april 2026)

---

## 2. Legacy Copilot Metrics API nedlagt

Det gamle Copilot Metrics API-et ble offisielt nedlagt 2. april 2026, som varslet i januar. Organisasjoner som fortsatt bruker de gamle endepunktene mister nå tilgang til bruksdata. Det nye Usage Metrics API-et leverer data via NDJSON-filer med langt mer detaljert telemetri — per språk, IDE, modell, kodelinje og redigeringsmodus.

Team-nivå-metrikker er ikke lenger tilgjengelig — kun organisasjons- og enterprise-nivå støttes i det nye skjemaet.

**Kilde:** [Closing down notice of legacy Copilot metrics APIs](https://github.blog/changelog/2026-01-29-closing-down-notice-of-legacy-copilot-metrics-apis/) (GitHub Changelog, 29. januar 2026)

---

## 3. Organisasjonsstyrt runner for cloud agent

Inntil nå ble runner-konfigurasjonen for coding agent satt per repository via `copilot-setup-steps.yml`. Nå kan organisasjonsadministratorer sette en standard-runner som brukes automatisk for alle repoer — og valgfritt låse innstillingen slik at individuelle repoer ikke kan overstyre den.

Dette gjør det enklere å rulle ut konsistente defaults (for eksempel større GitHub Actions-runnere for bedre ytelse) og sikre at agenten alltid kjører der organisasjonen vil — for eksempel på self-hosted runners med tilgang til interne ressurser.

**Kilde:** [Organization runner controls for Copilot cloud agent](https://github.blog/changelog/2026-04-03-organization-runner-controls-for-copilot-cloud-agent/) (GitHub Changelog, 3. april 2026)

---

## 4. Personvernpolicy for treningsdata trer i kraft

Den 24. april trer GitHubs oppdaterte personvernpolicy i kraft: interaksjonsdata fra Copilot Free-, Pro- og Pro+-brukere brukes til modelltrening med mindre de aktivt velger bort. Copilot Business og Enterprise er ikke berørt — kontraktsvilkårene beskytter enterprise-data.

Policyen ble kunngjort 25. mars og har møtt sterk kritikk for å være opt-out i stedet for opt-in. Utviklere som bruker personlige kontoer bør sjekke innstillingene under [Settings → Copilot → Privacy](https://github.com/settings/copilot).

**Kilde:** [Updates to GitHub Copilot interaction data usage policy](https://github.blog/news-insights/company-news/updates-to-github-copilot-interaction-data-usage-policy/) (GitHub Blog, 25. mars 2026)

---

## 5. Flere oppdateringer

- **Visual Studio Mars-oppdatering**: custom agents (`.agent.md`), agent skills, `find_symbol`-verktøy, Enterprise MCP governance med allowlist-policyer. [Kilde](https://github.blog/changelog/2026-04-02-github-copilot-in-visual-studio-march-update/)
- **GPT-5.1 Codex avviklet**: alle GPT-5.1-varianter (Codex, Codex-Max, Codex-Mini) er fjernet fra Copilot. Anbefalt erstatning er GPT-5.3-Codex. [Kilde](https://github.blog/changelog/2026-04-03-gpt-5-1-codex-gpt-5-1-codex-max-and-gpt-5-1-codex-mini-deprecated/)
- **Gemma 4 open source**: Google lanserer sin mest avanserte åpne modellfamilie under Apache 2.0 — fire varianter fra 2B til 31B parametere, multimodal (tekst, bilde, video, lyd), opptil 256K kontekst. [Kilde](https://blog.google/innovation-and-ai/technology/developers-tools/gemma-4/)

---

## Relevans for Nav

| Trend | Hva det betyr for Nav |
| --- | --- |
| Copilot SDK | Nav kan bygge egne verktøy med Copilots agentmotor — Go SDK er direkte relevant for mcp-onboarding og mcp-registry. Vurder for interne tjenester. |
| Legacy metrics API nedlagt | Navs copilot-metrics-app bruker allerede det nye Usage Metrics API-et — ingen handling nødvendig. Verifiser at ingen andre Nav-verktøy bruker det gamle API-et. |
| Org-runner for cloud agent | Sentralstyrt runner-konfigurasjon. Nav kan sette standard for alle repoer og låse til self-hosted runners ved behov — viktig for compliance og ytelse. |
| Personvernpolicy | Nav bruker Enterprise — ikke berørt. Informer utviklere som bruker personlige Copilot-kontoer om opt-out før 24. april. |
| GPT-5.1 deprecering | Sjekk om noen team har satt GPT-5.1 som foretrukket modell. |
