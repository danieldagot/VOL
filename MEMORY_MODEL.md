# VOL Memory Model (Prototype + Future Intent)

> Design document — **not** an ownership checker.
> Current runtime rules are **Supported** in [`SPEC.md`](SPEC.md).
> Future native/ownership answers are **Vision / Planned** until specified and tested.

Related: [`SPEC.md`](SPEC.md) §3.3 (arrays), §3.6 (structs), §3.7 (dicts), §5.2 (`const`),
[`IDEAS.md`](IDEAS.md) “Safety and Memory”.

---

## 1. What the interpreter does today

VOL’s tree-walker is dynamically typed. Values are Go `any` payloads. There is
**no** move checker, borrow checker, or garbage-collector contract exposed to
programs. Host GC (Go) reclaims unreachable values.

### 1.1 Scalars

Integers (`i64`), floats (`f64`), booleans, and strings behave as **immutable
values** under assignment: `b := a` copies the value. Rebinding `a` does not
change `b`. String indexing is rejected (`R003`).

### 1.2 Arrays

- Assignment (`:=` / `=`) and argument passing **share** the array reference.
- Indexed writes through any alias mutate the shared array.
- Rebinding a name to a different array updates only that binding.
- `.where` / `.map` return **new** arrays (eager).
- `.copy()` — shallow clone of the top-level array.
- `.deep_copy()` — recursive clone of nested arrays; nested **dicts** inside
  arrays are also cloned recursively (see SPEC §3.7).
- Move / ownership semantics are **out of scope** for the prototype.

### 1.3 Dicts

- Same **share-on-assign** rule as arrays.
- String-key get/set; missing key on get traps (`R045`).
- **No** dict `.copy()` / `.deep_copy()` today — gap vs arrays.
- Nested dicts are cloned when an enclosing array is `.deep_copy()`’d.

### 1.4 Structs

- Product structs are heap objects in the interpreter; assignment **shares**
  the instance. Field writes through aliases are visible to all aliases.
- `const` on a struct binding blocks rebinding the name, not field mutation.

### 1.5 `const` (shallow)

`const name := expression` prevents **reassignment** of `name` (`S030` / `R030`).
It does **not** freeze the underlying array, dict, or struct. Indexed / field
writes on a `const` binding remain allowed when the value is mutable.

### 1.6 Closures and escape

Anonymous / nested `fn` values capture the environment by reference (shared
bindings). Values that escape a scope remain reachable through the function
value. There is no stack-vs-heap distinction in the prototype.

### 1.7 Modules and `@std` handles

Opaque host handles (for example DB connections) are ordinary runtime values
passed by reference. Import bindings install module namespaces (SF-3.1); the
namespace object holds exported functions/values.

---

## 2. Design intent for a future native backend (Vision)

These are **not** implemented. They constrain future work so SF sugar does not
paint ownership into a corner.

| Question | Intent (Vision) |
| --- | --- |
| What owns an array/dict/struct? | A single owning binding or explicit shared handle; TBD spelling |
| What happens on assignment? | Prefer **move** of uniquely owned values; **share** only when annotated or inferred safe |
| Escape into closures? | Escaping values must be heap-allocated or copied; local non-escaping values may stay stack |
| Structs: value or reference? | Undecided; product structs today are reference-like — document before changing |
| Mutable aliases? | At most one mutable owner **or** explicit shared mutability; no silent UB |
| What does `const` mean under ownership? | Likely “immutable binding + deep immutability” or a separate `freeze` — must be redesigned vs today’s shallow `const` |
| Move vs copy? | Moves for unique ownership; `.copy` / `.deep_copy` remain explicit clones |
| Heap allocation? | Escape analysis first; explicit allocators later; **no** mandatory GC in the native story |

---

## 3. Gaps to close before claiming “systems memory safety”

1. Dict clone API parity with arrays (or a unified clone story).
2. Written ownership / aliasing rules (this file §2 → SPEC when decided).
3. Interaction of `const`, shared mutation, and future moves.
4. Closure capture rules under ownership.
5. FFI / host-handle lifetime.

Until then: describe VOL honestly as a **reference-sharing dynamically typed
prototype** with a documented path toward ownership — not as a safe systems
language runtime.

---

## 4. Non-goals for this document

- Implementing a borrow checker
- Changing SF-3.1 assignment-sharing behavior
- Defining GC or ARC as the native default
