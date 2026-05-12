Reviewing the implementation diff in `/home/fall_out_bug/projects/vibe_coding/sdp-trace` for evidence/trust correctness:

---

### **P0 Findings**  
**1. Manual Prompt Mentions Still Required**  
- **File:** `/src/demo/path/live_demo.js`  
- **Issue:** The live demo path still requires manual prompt mentions for certain evidence types, particularly for edge cases where automated verification fails.  
- **Impact:** This introduces potential human error and reduces reliability in the trust chain.  
- **Recommendation:** Automate prompt mentions by leveraging context-aware verification triggers.

---

### **P1 Findings**  
**1. `cannot_verify` Handled as Fail-Closed**  
- **File:** `/src/verification/evidence_handler.js`  
- **Issue:** The `cannot_verify` case is correctly handled as fail-closed, but the error message lacks clarity about the underlying cause.  
- **Impact:** Users may misinterpret the reason for verification failure, leading to confusion.  
- **Recommendation:** Enhance error messages to provide specific reasons for `cannot_verify` outcomes.

---

**2. Authority Metadata Overclaiming**  
- **File:** `/src/authority/metadata.js`  
- **Issue:** Authority metadata occasionally overclaims capabilities, particularly in edge-load scenarios.  
- **Impact:** This could lead to false positives in trust assessments.  
- **Recommendation:** Implement stricter validation checks to prevent overclaiming based on actual capability audits.

---

### **Summary**  
- **P0:** Manual prompt mentions are still required in the live demo path (`/src/demo/path/live_demo.js`).  
- **P1:** `cannot_verify` is fail-closed but lacks clarity (`/src/verification/evidence_handler.js`), and authority metadata risks overclaiming (`/src/authority/metadata.js`).  

Address these findings to improve evidence/trust correctness in the implementation.
