# API Gap Analysis Report
**Date**: 2026-03-29  
**Status**: Critical - 52 APIs non-functional

## Executive Summary

| Metric | Count |
|--------|-------|
| Backend Controllers | 103 files |
| Frontend Views | 279 pages |
| APIs Tested | 80 |
| **Working APIs** | **26** |
| **Broken APIs (404/500)** | **52+** |
| Coverage | ~33% |

---

## Working APIs (26) ✅

| Category | Endpoints |
|----------|-----------|
| Core | `/devices`, `/members`, `/stores`, `/departments` |
| OTA | `/ota/packages`, `/ota/deployments` |
| Alerts | `/alerts` |
| Dashboard | `/dashboard/stats` |
| System | `/system/health`, `/settings` |
| Users | `/users` |
| Auth | `/auth/login` (login works) |
| Organization | `/org/companies`, `/org/departments`, `/org/employees`, `/roles`, `/permissions` |
| Notifications | `/notification/channels`, `/notification/templates` |
| Content | `/knowledge`, `/knowledge/versions`, `/content/versions` |
| Subscriptions | `/subscriptions/auto-renewal/status` |
| Audit | `/audit/logs` |
| Announcements | `/announcements` |
| Certificates | `/certificates` |
| AI | `/ai/config` |
| Misc | `/position-templates` |

---

## Broken APIs (52+) ❌

### Priority 1 - Core Business Logic (12)
| API | Issue | Likely Cause |
|-----|-------|--------------|
| `/ai/chat` | 404 | Route not registered or handler nil |
| `/ai/models` | 404 | Route not registered |
| `POST /stores` | 500 | Database error |
| `POST /members` | 500 | Database error |
| `POST /users` | 500 | Database error |
| `/org/employees` | 500 | Database error |
| `/org/posts` | 500 | Database error |
| `/pet` | 404 | Route mismatch (pets vs pet) |
| `/pet/profile` | 404 | Route not registered |
| `/action/library` | 404 | Route not registered |
| `/action/learning` | 404 | Route not registered |
| `/simulation` | 404 | Route not registered |

### Priority 2 - Enterprise Features (18)
| API | Issue |
|-----|-------|
| `/digital-twin/*` | 404 (despite route existing) |
| `/emotion/*` | 404 |
| `/voice-emotion` | 404 |
| `/health` | 404 |
| `/simulation` | 404 |
| `/map` | 404 |
| `/insurance` | 404 |
| `/analytics` | 404 |
| `/reports` | 404 |
| `/advanced` | 404 |
| `/batch` | 404 |
| `/flow` | 404 |
| `/schedule` | 404 |
| `/regions` | 404 |
| `/timezones` | 404 |
| `/data-residency` | 404 |
| `/integration` | 404 |
| `/subscription/gift` | 404 |

### Priority 3 - Device Enhancements (8)
| API | Issue |
|-----|-------|
| `/device/health` | 404 |
| `/device/security` | 404 |
| `/device/monitor` | 404 |
| `/device/shadow` | 404 |
| `/ota/compatibility` | 404 |
| `/ota/sdk` | 404 |
| `/security/evo` | 404 |
| `/platform/evo` | 404 |

### Priority 4 - Pet & Family (6)
| API | Issue |
|-----|-------|
| `/pet/social` | 404 |
| `/family/album` | 404 |
| `/child/mode` | 404 |
| `/smart/home` | 404 |
| `/miniclaw` | 404 |
| `/market` | 404 |
| `/miniapp` | 404 |

### Priority 5 - Compliance & Security (8)
| API | Issue |
|-----|-------|
| `/compliance` | 404 |
| `/policy` | 404 |
| `/gdpr` | 404 |
| `/data-masking` | 404 (route exists) |
| `/alerts/dedup` | 404 |
| `/alerts/healing` | 404 |
| `/alert/rules` | 404 |
| `/alert/history` | 404 |
| `/alert/settings` | 404 |

### Priority 6 - Member Enhancements (6)
| API | Issue |
|-----|-------|
| `/member/profile` | 404 |
| `/member/enhanced` | 404 |
| `/card` | 404 |
| `/coupon` | 404 |
| `/content/review` | 404 |
| `/offline/sync` | 404 |

---

## Root Cause Analysis

### 1. Route Registration Issues
Many controllers exist but their `RegisterRoutes()` methods are never called in `main.go`:
- `EmotionController` - registered but routes not called
- `DigitalTwinController` - registered but routes not called
- `HealthController` - registered but routes not called
- `SimulationController` - registered but routes not called

### 2. Path Mismatch
- Frontend calls `/api/v1/ai/chat`
- Backend route pattern might be `/ai/chat` (without `/api/v1` prefix at route level)

### 3. 500 Errors - Database Issues
- `/stores`, `/members`, `/users` POST return 500
- Likely: missing tables, column mismatches, or constraint violations

### 4. Handler Panics/Nil
- Some routes appear registered but return 404
- Likely: controller handler function is nil or panics

---

## Recommended Fix Sequence

### Phase 1: Fix 500 Errors (Blocking)
1. `POST /stores` - Check `stores` table structure
2. `POST /members` - Check `members` table structure  
3. `POST /users` - Check `users` table structure
4. `/org/employees` - Check `employees` table
5. `/org/posts` - Check `posts` table

### Phase 2: Register Missing Routes
1. AI: `/ai/chat`, `/ai/models`
2. Digital Twin: `/digital-twin/*`
3. Emotion: `/emotion/*`
4. Health: `/health/*`
5. Simulation: `/simulation/*`
6. Pet: `/pets`, `/pets/profile`, `/pet/social`
7. Action: `/action/library`, `/action/learning`

### Phase 3: Fix Path Mismatches
1. Check all route registrations match API contract
2. Ensure `/api/v1` prefix consistency

### Phase 4: Frontend Validation
1. Test all frontend API calls match backend routes
2. Fix any mismatched endpoints

---

## Database Tables Status
- Total: 319 tables (per docs)
- All tables should exist per migration files
- 500 errors suggest data model vs code mismatch

---

## Files to Modify

### Priority Fix Files (main.go + device_controller.go)
1. Add route registrations for missing controllers
2. Fix 500 errors by checking DB queries
3. Verify handler functions are not nil

### Controllers Needing Route Registration
```
action_learning_controller.go
action_library_controller.go
digital_twin_controller.go
emotion_controller.go
health_controller.go
voice_emotion_controller.go
simulation_controller.go
pet_controller.go
pet_profile_controller.go
pet_social_controller.go
device_health_score_controller.go
device_security_controller.go
device_monitor_controller.go
device_shadow_snapshot_controller.go
insurance_controller.go
offline_controller.go
...
```

---

_Report generated: 2026-03-29_
