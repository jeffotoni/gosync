# gosync

A principio estamos desenvolvendo um protótipo, e simulando o comportamento de um sync para buckets na nuvem.
Iremos utilizar neste projeto o S3 da AWS e o Space da DigitalOcean, o nosso programa irá enviar para nuvem arquivos e sicroniza-los do seu Local para Nuvem e vice versa.

Lembrando que este processo não é tão trivial como aparentemente apresenta ser.

Temos alguns desafios que temos que pontuar para que possamos desenvolver.

O objetivo aqui é evoluir o desenvolvimento a medida que vamos descobrindo as melhores formas de implementar em Go o sync.