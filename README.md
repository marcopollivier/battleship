# Battleship

<img src="doc/img/gopher_battleship.png" width="30%">

Até os anos 90 era comum encontrar em papelarias um bloco de papel que continha instruções para um jogo simples, mas que poderia vir a se tornar bem divertido se jogado com as pessoas certas. Nas folhas desse bloco continha duas matrizes de `10x10` organizadas entre numeros e letras, ou seja, matrizes de `10xJ`.

Uma dessas matrizes representava as casas onde você dispõe suas "peças" e nas outras você faz marca as jogas que você faz em seu turno. Além das duas matrizes, o bloco possuia uma indicação de como deveria ser o modelo de cada uma das peças que você deveria utilizar. Uma legenda por assim dizer. As formas das peças variavam e os próprios nomes também. Porta-aviões, destroyers, contratorpedeiros, rebocadores, submarinos etc. Todos poderiam participar dependendo da versão que havia sido definida.

Pronto. Nisso se compunha um jogo de **Batalha Naval** antes de você ter opções sofisticadas que hoje rodam até mesmo em smartphones de formas muito mais elaboradas que um simples bloco de papel.

No ano de 1990 e pouco, esse era o jogo que unia a mim e ao meu pai durante dias. Meu pai era professor na extinta Escola Politécnica de Araruama e também era um programador rudimentar de [Basic](https://pt.wikipedia.org/wiki/BASIC), mas que fazia bem o que se propunha a fazer. E em uma dessas oportunidades, ele fez um Jogo de Batalha Naval em Basic - com interface gráfica e tudo. Infelizmente na época não existia um repositório como esse para guardar seus códigos e o jogo foi gravado em um frágil disquete 3,5 e se perdeu pouco depois de ter sido construído. Uma tragédia lamentável. Como eu não tinha computador em casa, ficava limitado a jogar no laboratório da escola e, portanto, não aproveitei tantas vezes quanto gostaria.

Bom. Memórias boas não precisam ficar apenas nas nossas mentes. Então o propósito desse projeto é recriar um jogo simples de Batalha Naval para reavivar a memória de um pai tão magnífico quanto o meu e, de quebra, me fazer voltar alguns anos no tempo em que isso era uma atividade que tanto nos unia. A diferença é que esse projeto não será feito em Basic, mas em [Go](https://go.dev) (desculpa, pai, mas Basic já é demais). E ainda assim, acredito que será um bom paralelo, pois Basic era a linguagem que ele se dedicava e, 20 anos depois, Go é a que eu tenho me dedicado.

Em memória de Julio Ollivier.

## O jogo

O tabuleiro do jogo é constituído por uma matriz de 10x10 onde as escalas são dispostas da seguinte forma:

- Colunas verticais marcadas de 1 a 10
- Colunas horizontais de A a J

Este jogo é recomendado para duas pessoas. Ambos os jogadores devem posicionar seus Navios de Combate em sua matriz 10xJ.

- O primeiro jogador dispara indicando uma posição da Matriz. Exemplo: A1, B3, J10.
- Caso ele acerte, tem o direito de efetuar mais um disparo.
- Caso erre, é considerado um tiro na água e a vez passa a ser do outro jogador.
- Quando todas as partes de um Navio são acertadas o mesmo é considerado como "Afundado"
- Quem afundar todos os Navios do oponente, **ganha o jogo**

A frota segue a cartela de papel original:

<img src="doc/img/batalha-naval.jpg" width="60%">

Navio            | Tamanho | Quantidade | Formato
---------------- | ------- | ---------- | -------
Submarino        | 1       | 3          | 1 casa
Destroyer        | 2       | 2          | Linha de 2 casas
Cruzador         | 3       | 1          | Linha de 3 casas
Encouraçado      | 4       | 1          | Linha de 4 casas
Porta-Aviões     | 5       | 1          | T: 3 casas em linha e 2 descendo do meio

São 8 navios, que ocupam 19 casas. Os navios em linha podem ser posicionados na horizontal ou na vertical. O Porta-Aviões pode ser posicionado em qualquer uma das 4 orientações do T:

```
■■■     ■       ■   ■
 ■      ■     ■■■   ■■■
 ■     ■■■      ■   ■
```

### Submarino

**Descrição**: Submarino é uma embarcação especializada para operar submersa, tendo sido largamente usada pela primeira vez na Primeira Guerra Mundial, sendo empregada por todas as grandes marinhas atualmente. Submarinos civis e submergíveis são usados com fins científicos tanto na água doce quanto salgada para trabalhar em profundidades muito grandes para mergulhadores humanos. **Fonte:** [Wikipedia](https://pt.wikipedia.org/wiki/Submarino)

### Destroyer

**Descrição**: Um contratorpedeiro, também referido como Destroyer ou Destróier, é um tipo de navio de guerra, rápido e manobrável, concebido para escoltar navios maiores numa esquadra naval ou comboio de navios e defendê-los contra agressores menores, mas perigosos. **Fonte:** [Wikipedia](https://pt.wikipedia.org/wiki/Contratorpedeiro)

AKA: Contratorpedeiro, Destróier

### Cruzador

**Descrição**: Um cruzador é um tipo de navio de guerra. Designando originalmente qualquer navio encarregado de tarefas de exploração em uma esquadra, o termo "cruzador" passou a designar, no final do século XIX, um tipo específico de navio oceânico, maior e mais armado que as fragatas, que acabou por substituí-las. **Fonte:** [Wikipedia](https://pt.wikipedia.org/wiki/Cruzador)

### Encouraçado

**Descrição**: Um couraçado ou encouraçado é um navio de guerra pesadamente blindado e armado com as peças de artilharia de longo alcance e de maior calibre existentes. Normalmente, os couraçados eram maiores, mais armados e mais blindados que os cruzadores e contratorpedeiros. **Fonte:** [Wikipedia](https://pt.wikipedia.org/wiki/Coura%C3%A7ado)

AKA: Couraçado

### Porta Aviões

**Descrição**: Um porta-aviões é um navio de guerra cujo papel principal é servir de base aérea móvel. Permite, portanto, que uma força naval possa projetar o seu poderio aéreo a grandes distâncias, sem a necessidade de depender de aeroportos para os seus aviões. **Fonte:** [Wikipedia](https://pt.wikipedia.org/wiki/Porta-avi%C3%B5es)

## TODO - Evolução do projeto

- Permitir que uma pessoa ataque um tabuleiro armado aleatoriamente
- Permirir que arme seu próprio tabuleiro e jogue contra a máquina e jogadas aleatórias
- Permitir que que a pessoa definam a quantidade de navios que ela quer jogar
