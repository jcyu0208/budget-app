# Feature workflow (read before starting any feature)

docs/specs.md is the high-level baseline. It is a starting point, not the full requirements. Assume it has gaps or errors.

1. Interview. Read the spec section for this feature. Look for gaps, ambiguities and flaws in what it claims, for this feature and for adjacent features and past decisions. Ask me questions one at a time. Never fill a gap silently. For each real choice, give 2 or 3 options with tradeoffs and a recommendation. Cite the YNAB help page for any YNAB rule, and say if you are working from memory. If a page won't load within a few minutes, say so and move on. Do not crawl.
2. Plan. Use Plan Mode. No code until I approve. Show the approach, which core rules, ports and adapters change, and the test plan.
3. Test plan. Go beyond unit tests.
   - Functional tests drive the real entry point (the API) against the in-memory adapter, using the worked example numbers. They assert through the read side. After a transaction, query the Available balance. Do not inspect stored values.
   - After every step in a functional test, check the invariant: accounts = Ready to Assign + all Available.
   - Contract tests run the same suite against every storage adapter.
   - Cover the spec's edge cases: zero, negative, empty, duplicate, very large.
4. Architecture. Clean architecture with ports and adapters. The core holds the rules and defines the ports. It imports no HTTP and no SQL. Storage is in memory first, with SQL later as an adapter. Swapping storage must need no core change. Numbers are derived from records, never stored. Money is whole cents. Months exist from the first ticket. Edits and deletes are out of scope until I say so.
5. Build. Test first, on a feature branch, in small commits.
6. Feature doc. Write docs/features/<name>.md before building and update it after: behaviour, rules, edge cases, decisions, test plan, links to code and tests. One page at most, half a page for a ticket under 45 minutes.
7. Open items. Log every parked idea, question and assumption in docs/open-items.md as a row: id, date, type (question, assumption or parked idea), feature, text, status. Log every assumption you make while building. At the start of each feature, review open items and close or carry them.
