# jankenpo
Jankenpo game to demonstrate the use of different approaches in distributed systems

O objetivo do projeto é criar um jogo de "Pedra, papel e tesoura", onde o cliente pede 2 valores, 
que são enviados para o servidor. 
O servidor então processa as jogadas e retorna: quem foi o ganhador (o player 1 ou o player 2), 
se foi empate, ou se teve alguma jogada inválida.

Obs.: cada pasta trata de um tipo de conexão diferente e tem seus próprios arquivos de servidor e cliente.
Obs.2: Pode ser executado de forma automatizada ou manual, basta alterar o parâmetro "auto" em PlayJanKenPo
para true ou false respectivamente
