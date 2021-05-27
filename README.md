# COVID-19-Simulation

Projekat iz predmeta Napredne tehnike programiranja.
</br>
<h2>Opis aplikacije</h2>
COVID-19 Simulation predstavlja aplikaciju koja simulira širenje zaraze virusom COVID-19. Na početku može postojati jedan nulti pacijent, a može ih biti i više. Neke od jedinki se mogu kretati, dok druge miruju. Na mapi postoji centar okupljanja(škola, tržni centar, ...) gde jedinke često dolaze. Kretanje jedinki je takvo da u svakoj iteraciji jedinka ima najveću šansu da se uputi ka centru okupljanja, dok su šanse za odlazak na ostala mesta na mapi manje. </br> </br>Karakteristike epidemije su sledeće: </br>
• duration - broj dana koliko će jedinka biti bolesna ukoliko se zarazi, nakon što protekne zadati broj dana, jedinka se oporavlja ili umire </br>
• incubation - broj dana koliko se jedinka nalazi u inkubaciji </br>
• infection rate - verovatnoća da će zaražena jedinka zaraziti drugu jedinku koja se nalazi u njenoj okolini </br>
• mortality - verovatnoća da jedinka umire nakon zaraze </br>
• immunity - verovatnoća da se zarazi jedinka koja je prethodno već bila zaražena</br>  </br>
U simulaciji, postoji i bolnica koja je određenog kapaciteta. Tokom svakog dana (iteracije) testiraju se random jedinke. Verovatnoća da se otkrije da zaražena jedinka ima COVID19 je velika, u ovom slučaju, ukoliko postoje slobodna mesta u bolnici, jedinka biva hospitalizovana, i njene šanse da preživi se povećavaju na 90%. U slučaju da nema mesta u bolnici, jedinka se šalje u karantin i tada joj je onemogućeno kretanje, ali njene šanse da preživi nisu uvećane. U slučaju da je doneta odluka o socijalnom distanciranju, šanse da jedinke koje se kreću idu ka centru okupljanja bivaju male.
</br>

<h2>Arhitektura sistema</h2>
Sistem bi se sastojao iz dva dela: logičkog dela koji je implementiran uz pomoć Golang programskog jezika, dok će sama vizualizacija simulacije biti predstavljena pomoću Pharo. </br>
<h3>Go logička jedinica</h3>
Logička jedinica bi se sastojala od jednog Go procesa.
Glavni zadaci logičke jedinice u toku iteracije su kretanje jedinki, donošenje odluka i obrada podataka o samoj jedinki na osnovu karakteristika epidemije.
</br> 
<h3>Pharo deo</h3>
Pharo deo podrazumeva vizualizaciju epidemije kroz iteracije. Takođe, obuhvata i analitičke prikaze, kao što su broj obolelih po danima, broj umrlih po danima itd.
