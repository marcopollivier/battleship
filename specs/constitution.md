# Constituição do projeto

Princípios que todo plano precisa respeitar. Um plano que viole algum deles precisa justificar a exceção na seção "Checagem da constituição" e ter aprovação explícita.

1. **Domínio separado da interface.** Regras do jogo (tabuleiro, navios, disparos, turnos) vivem em pacotes puros, sem `fmt.Print`, leitura de stdin ou aleatoriedade global. A interface de terminal só chama o domínio e exibe o resultado.
2. **Determinismo testável.** Tudo que for aleatório (posicionamento, jogadas da máquina) recebe a fonte de aleatoriedade por parâmetro, para que os testes usem uma seed fixa.
3. **Teste primeiro.** Cada critério de aceitação tem um teste escrito antes do código. Preferir testes table-driven com `testing` da biblioteca padrão.
4. **Somente biblioteca padrão.** Nenhuma dependência externa sem justificativa no plano.
5. **Erros como valores.** Entradas inválidas (coordenada fora do tabuleiro, navio sobreposto, tiro repetido) retornam `error`, nunca `panic`.
6. **Incrementos pequenos.** Uma spec entrega um comportamento completo e testado. Se o plano tiver mais de ~8 tarefas, a spec deve ser dividida.
7. **Idiomas.** Código, identificadores e comentários em inglês (o lint usa `misspell` com locale US). Specs, documentação e textos exibidos ao jogador em português.
8. **Lint limpo.** `make pre-commit` passa antes de cada commit.
