# Scraper do pobierania arkuszy maturalnych z Arkusze.pl
Pobiera wszystkie arkusze maturalne w formacie PDF z Arkusze.pl, normalizuje nazwy plików i organizuje je w odpowiednich folderach.

> [!NOTE]  
> Pobiera tylko arkusze maturalne. Aby pobrać inne arkusze, trzeba zmodyfikować lekko kod.

## Wymagania
- golang

## Uruchomienie
```shell
go run .
```
Pliki znajdą się w folderze `arkusze/` w aktualnej ścieżce.

## Struktura plików

```
arkusze/
|--przedmiot/
    |--pytania/ - główne arkusze
    |--odpowiedzi/ - karty z odpowiedziami
    |--transkrypcje/ - transkrypcje dla języków obcych
    |--dodatki/ - informatory, tablice, wzory, mapy
```

## Nazewnictwo plików

[rok]-[miesiąc]-[atrybuty (np. rozszerzona, poprawa...)].pdf


## Licencja
MIT (C) 2024 Maximilian Gaedig
