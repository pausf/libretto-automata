# close-flow-open-questions

Tracker: none

## The ask, verbatim

> haz toda esta tareas en solo una pr
>
> 1. Cerrar la divergencia del drift con un opt-in, no con un default. spec-drift avisa
> y nunca bloquea, y la razón (un gate que sorprende se desinstala) es buena. Pero hoy
> el que quiere que bloquee tiene que cablearlo solo. Yo shippearía el snippet listo —
> una línea documentada de hook pre-commit en el README del skill, copy-paste. El
> default no cambia; la distancia entre "quiero enforcement" y tenerlo baja a cero.
> Costo: unas líneas de docs.
>
> 2. Resolver el open question que ya tenés escrito: dónde se mira el artefacto.
> FLOW.md lo confiesa al final — la primera paleta cumplía su spec y era ilegible,
> 1.4:1 de contraste, y lo que la cazó fue un script descartable, no el flow. El
> reviewer del seam lee specs y diffs, no píxeles. Es el único agujero real que el
> flow se conoce a sí mismo y sigue abierto. Yo lo capturaría en la cola hoy: una
> regla en fase 6 — "si el cambio tiene salida visual, renderizala y mirala antes del
> seam" — no una fase nueva.
>
> 3. Usar metrics para decidir la pregunta que trajiste vos. Querías más comunicación
> en spec y plan. Antes de relitigar la decisión de 2026-08, medí: ¿dónde aparecen las
> correcciones — en las tres preguntas de fase 2, en findings del reviewer, o en
> cambios tuyos post-push? Si los findings tardíos son cosas que una pregunta de fase
> 2 hubiera evitado, eso es la evidencia para subir de tres preguntas. Si no, la
> decisión queda. Ya construiste el instrumento; apuntalo a esto. y se mas humano en
> las specs es decir no te de miedo preguntar mas, ya que si no eres demasiado
> automata en la primeras fases, recuerda que las preguntas siempre con claude
> question

## The reading

Three tasks, one PR, none of them a bug — all three close gaps the flow has already
declared about itself:

1. **Ship the opt-in drift gate as docs.** A documented copy-paste pre-commit hook
   snippet in `skills/record-work/`, so whoever wants `spec-drift` to block can have it
   in one paste. The warn-never-block default does not move.
2. **Answer FLOW.md's open question.** A rule inside phase 6 — a change with visual
   output gets rendered and looked at before the review seam — not a new phase. Closes
   the "where does the artifact get looked at?" item in FLOW.md's Open section.
3. **Point `libretto metrics` at correction attribution**, so the ask-more-questions
   debate is decided by evidence: where do corrections surface — phase 2 answers,
   reviewer findings, or post-push commits? And loosen phase 2's fear of asking: the
   contract is built by two people, questions always via AskUserQuestion.
