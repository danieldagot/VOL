# Example roast (tone and structure reference)

The following is a gold-standard example of voice, structure, and severity.
When reviewing, match this energy against **current** docs — do not copy outdated
findings if the repo has fixed them.

**Current doc layout (use this, not the filenames inside the sample body):**
`README.md`, `SPEC.md`, `IDEAS.md`, `AGENTS.md`, plus spot-checks of
`internal/lang/` and `examples/`. `SYNTAX.md` and `VOCABULARY.md` were removed;
implemented grammar and vocabulary live in `SPEC.md` only.

The sample below is historical tone reference. Some concrete bugs it names may
already be fixed; re-verify everything.

---

I read the current README, SYNTAX.md, VOCABULARY.md, IDEAS.md, and AGENTS.md.

Bottom line: the idea behind VOL is much better than the current language specification. Right now, VOL is a strong language manifesto attached to a small interpreter, not yet a serious systems-language specification. The repo itself accurately admits that the static type system, ownership model, explicit types, errors, optionals, structs, generics, async, allocator control, and backend are still open.

My score today
Area	Score
Core idea	8/10
Surface syntax	7/10
Simplicity	8/10
Novelty	5/10
Specification quality	3/10
Systems-language semantics	2/10
Actual evidence of LLM optimization	1/10
Potential	9/10

The harshest version of the review is:

VOL is currently a pitch deck wearing a parser.

And that's completely fixable at this stage.

The roast
You don't actually have a language specification yet. SYNTAX.md says that operator precedence “must be formally defined,” mutability isn't decided, if-as-expression isn't decided, while isn't final, parameter/return types aren't final, and major constructs are undecided. There's also no formal lexical grammar or EBNF. For something intended to be especially deterministic for LLMs, that is a major hole. An LLM should be able to answer exactly whether a + b * c, a +\n b, foo\n.bar(), or not a == b is legal and what it means.

“Token optimized for LLMs” is currently an assertion, not a property. AGENTS.md says “Every token matters,” proposes build-time token counts and a future “Meaning per Token” metric—but MPT isn't defined. This is possibly the biggest conceptual weakness in VOL. Token counts depend on the tokenizer. GPT, Claude, Gemini, Llama, etc. won't necessarily tokenize :=, .where, each, &&, and and identically. More importantly, shortest source != best LLM language. If VOL uses 15% fewer tokens but causes 30% more compiler errors, it has failed. Your real metric should look more like:

task success / total tokens consumed, including generated code, compiler errors, repair prompts, and revisions.

“Compiler infers ownership, borrowing and lifetimes” is doing absurd amounts of work. The vision says the compiler handles ownership analysis, borrow analysis, lifetime inference, escape analysis, vectorization and parallelization while avoiding a GC and retaining C-like performance. That's basically several PhDs hiding behind the word “infer.” Local ownership inference? Plausible. Eliminating ownership contracts from public APIs? Much harder. Eventually the specification must answer: who owns returned memory, can references escape, can two mutable aliases exist, what happens when inference is ambiguous, whether something gets copied, whether RC is permitted, and how FFI works. Those are language semantics, not optimizer implementation details.
Your nicest feature, where, is also a semantic landmine. I like this:
total := numbers.where(_ > 5).sum

It's much closer to what VOL should be than conventional loops. But what does where return? New array? Lazy view? Iterator? Compiler IR that gets fused with sum? Does it allocate? Can the predicate mutate external state? What happens with nested _? Is order guaranteed? Can the compiler parallelize it? What happens on integer overflow during sum?

For a normal scripting language, you can hand-wave some of that. For a language promising predictable native systems performance and no hidden unbounded allocation, you absolutely cannot.

“Vector-Oriented Language” currently isn't vector-oriented. The README calls it Vector-Oriented Language, while AGENTS.md calls it Vocabulary Optimized Language, and the README also jokes that it's Vibe-Oriented Language. Pick one identity. More importantly, .each, .where, and .sum are collection operations, not a vector programming model. If “Vector” means the compiler recognizes bulk operations and performs SIMD/vectorization, make that an actual semantic pillar. Otherwise the name promises something the language doesn't currently contain.
The type system is the elephant missing from the room. Currently your literal model is essentially i64, f64, bool, string, arrays. A systems language eventually needs things like fixed-width integers, signed/unsigned behavior, conversion rules, layouts, structs, enums, references, pointers/handles, slices/views, alignment and ABI rules. Even the innocent-looking statement that string .length counts “characters” is under-specified. Does “character” mean bytes, Unicode scalar values, or grapheme clusters? In a systems language, that distinction is enormous.
Your “one canonical way” principle already contradicts the language. You explicitly say to prefer one canonical way of expressing an operation. Then you allow both:
numbers.each number {
    ...
}

and

numbers.where(_ > 5).sum

and tell programmers to choose whichever communicates intent better. That's not necessarily bad language design—but then “one canonical way” isn't your principle. I'd change the principle to “one canonical representation for each distinct intent.” Filtering/summing and imperative iteration are genuinely different intents.

Your own documentation rules are already being violated. AGENTS.md says planned behavior belongs in IDEAS.md, while SYNTAX.md should contain implemented/provisional grammar; “Provisional” specifically means implemented behavior whose spelling or semantics may change. Yet SYNTAX.md contains a parallel { ... } syntax whose semantics haven't been decided, while VOCABULARY.md explicitly says parallel is only a design idea and isn't accepted by the interpreter. That's exactly the kind of ambiguity an AI-first language should aggressively eliminate.
And there is an actual hilarious contradiction in VOCABULARY.md. The example labels Double as a “public function” because it's capitalized, then immediately says capitalization has no visibility meaning and functions are public only through export. Also the or row in the vocabulary table is mangled because the C-style || clashes with the Markdown table separators. Your AI-optimized language specification has been defeated by two pipe characters. 😄
The biggest design issue

I would stop adding language features right now.

Before parallel, async, generics, macros, HTTP syntax, or anything exciting, make VOL v0.1 semantically boring and ridiculously precise.

Take just this:

values
variables
functions
if
repeat
while
arrays
each
where
sum

And specify every corner of it.

In particular, create a real SPEC.md containing lexical grammar → formal expression grammar → type/value semantics → scopes → evaluation order → numeric behavior → arrays/string semantics → mutation → function semantics → control flow → failure behavior.

Then define one second document that may be even more important for VOL:

LLM_BENCHMARK.md

Your language's killer contribution shouldn't merely be unusual syntax. It should be the first language where you can say:

“On these 100 programming tasks, VOL required 34% fewer generation tokens than Go, 41% fewer than Rust, produced 22% fewer compile failures, and required 37% fewer repair tokens.”

That would make VOL interesting.

Without that benchmark, “LLM optimized” is marketing.

With it, VOL becomes an experiment people can actually falsify, measure and improve.

And one final opinion: don't try too hard to make the syntax look alien. :=, braces, fn, if, return, arrays, and normal expressions are good precisely because models have billions of tokens of training experience with similar constructs. Your innovation should be semantic compression like .where(...).sum, compiler inference, structured diagnostics, canonicalization, and higher-level intent—not inventing a weird glyph for everything.

So I would keep about 80% of your existing surface syntax.

I would throw out about 80% of the claims that the compiler will magically infer everything until you've written the semantic rules proving how it can.

That's where I think VOL could become genuinely interesting.
