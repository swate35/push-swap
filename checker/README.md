# Checker

`checker.go` est le vérificateur du projet **push-swap**. Il prend les nombres
à trier en arguments, lit une suite d'instructions sur l'entrée standard et
simule ces instructions sur deux piles (`A` et `B`).

Depuis la racine du dépôt, lancez le programme avec Go :

```powershell
go run .\checker\checker.go "3 2 1"
```

Les nombres peuvent aussi être fournis comme arguments séparés :

```powershell
go run .\checker\checker.go 3 2 1
```

Dans les deux formes, le premier nombre est placé au sommet de la pile `A`.

## Envoyer les instructions

Les instructions doivent être saisies une par ligne sur l'entrée standard.
Dans PowerShell, cette commande attend une saisie interactive :

```powershell
go run .\checker\checker.go "3 2 1"
```

Après avoir lancé la commande, saisissez les instructions une par ligne.
Pour terminer la saisie dans le terminal Windows, appuyez sur `Ctrl+Z`,
puis sur Entrée.

Pour fournir directement plusieurs instructions, utilisez les éléments d'un
tableau PowerShell :

```powershell
@("sa", "rra") | go run .\checker\checker.go "3 2 1"
```

Exemple avec une suite permettant de trier `3 2 1 0` :

```powershell
@("rra", "pb", "sa", "rra", "pa") | go run .\checker\checker.go "3 2 1 0"
```

Sortie attendue :

```text
OK
```

Un petit script PowerShell peut aussi écrire les instructions une par une :

```powershell
@("sa", "rra", "pb") | go run .\checker\checker.go "3 2 1"
```

## Vérifier la sortie de push-swap

Depuis la racine du dépôt, une fois le programme `push-swap` compilé, ses
instructions peuvent être envoyées directement au checker :

```powershell
$ARG = "4 67 3 87 23"
.\push-swap.exe $ARG | go run .\checker\checker.go $ARG
```

Le checker affiche alors `OK` si le résultat de `push-swap` trie correctement
la pile.

## Compiler le checker

Pour produire un exécutable nommé `checker.exe` :

```powershell
go build -o checker.exe .\checker\checker.go
```

Puis lancez-le :

```powershell
@("rra", "pb", "sa", "rra", "pa") | .\checker.exe "3 2 1 0"
```

## Instructions reconnues

| Instruction | Effet |
| --- | --- |
| `sa`, `sb` | Échange les deux premiers éléments de A ou B. |
| `ss` | Effectue `sa` et `sb`. |
| `pa`, `pb` | Déplace le sommet de B vers A, ou de A vers B. |
| `ra`, `rb` | Décale le sommet de A ou B vers le fond. |
| `rr` | Effectue `ra` et `rb`. |
| `rra`, `rrb` | Décale le fond de A ou B vers le sommet. |
| `rrr` | Effectue `rra` et `rrb`. |

Une opération de rotation ou d'échange n'a aucun effet si la pile concernée
contient moins de deux éléments. Un `push` sur une pile source vide n'a aucun
effet.

## Résultat

- `OK` si A est triée par ordre croissant et B est vide à la fin de la
  simulation.
- `KO` dans les autres cas.
- `Error` (sur la sortie d'erreur) si un argument est invalide, si une valeur
  est dupliquée, si une instruction est inconnue ou si la lecture de l'entrée
  échoue.

Les valeurs doivent être des entiers signés sur 32 bits. Sans argument, le
programme se termine sans afficher de résultat.
