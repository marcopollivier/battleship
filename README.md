# Battleship

Até os anos 90 era comum encontrar em papelarias um bloco de papel que continha instruções para um jogo simples, mas que poderia vir a se tornar bem divertido se jogado com as pessoas certas. Nas folhas desse bloco continha duas matrizes de `10x10` organizadas entre numeros e letras, ou seja, matrizes de `10xJ`. 

Uma dessas matrizes representava as casas onde você dispõe suas "peças" e nas outras você faz marca as jogas que você faz em seu turno. Além das duas matrizes, o bloco possuia uma indicação de como deveria ser o modelo de cada uma das peças que você deveria utilizar. Uma legenda por assim dizer. As formas das peças variavam e os próprios nomes também. Porta-aviões, destroyers, contratorpedeiros, rebocadores, submarinos etc. Todos poderiam participar dependendo da versão que havia sido definida. 

Pronto. Nisso se compunha um jogo de *Batalha Naval* antes de você ter opções sofisticadas que hoje rodam até mesmo em smartphones de formas muito mais elaboradas que um simples bloco de papel. 

No ano de 1990 e pouco, esse era o jogo que unia a mim e ao meu pai durante dias. Meu pai era professor na extinta Escola Politécnica de Araruama e também era um programador rudimentar de [Basic](https://pt.wikipedia.org/wiki/BASIC), mas que fazia bem o que se propunha a fazer. E em uma dessas oportunidades, ele fez um Jogo de Batalha Naval em Basic - com interface gráfica e tudo. Infelizmente na época não existia um repositório como esse para guardar seus códigos e o jogo foi gravado em um frágil disquete 3,5 e se perdeu pouco depois de ter sido construído. Uma tragédia lamentável. Como eu não tinha computador em casa, ficava limitado a jogar no laboratório da escola e, portanto, não aproveitei tantas vezes quanto gostaria. 

Bom. Memórias boas não precisam ficar apenas nas nossas mentes. Então o propósito desse projeto é recriar um jogo simples de Batalha Naval para reavivar a memória de um pai tão magnífico quanto o meu e, de quebra, me fazer voltar alguns anos no tempo em que isso era uma atividade que tanto nos unia. A diferença é que esse projeto não será feito em Basic, mas em Go (desculpa, pai, mas Basica já é demais). E ainda assim, acredito que será um bom paralelo, pois Basic era a linguagem que ele se dedicava e, 20 anos depois, Go é a que eu tenho me dedicado.

Em memória de Julio Ollivier.

## O jogo 
O tabuleiro do jogo é constituído por uma matriz de 10x10 onde as escalas são dispostas da seguinte forma: 
- Colunas verticais marcadas de 1 a 10 
- Colunas horizontais de A a J 



Os navios são posicionados de acordo com seu tamanho e posição no tabuleiro tendo os seguintes tamanhos e quantidade inicial:

Navio            | Tamanho | Quantidade
---------------- | ------- | ----------
Submarino        | 1       | 4
Destróier        | 2       | 3
Contratorpedeiro | 3       | 2
Encouraçado      | 4       | 2
Porta-Aviões     | 5       | 1


## TODO - Evolução do projeto

- Permitir que uma pessoa ataque um tabuleiro armado aleatoriamente 
- Permirir que arme seu próprio tabuleiro e jogue contra a máquina e jogadas aleatórias
- Permitir que que a pessoa definam a quantidade de navios que ela quer jogar 
