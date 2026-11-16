# Checker

`checker.go` est le vérificateur du projet **push-swap**. Il prend les nombres
à trier en arguments, lit une suite d'instructions sur l'entrée standard et
simule ces instructions sur deux piles (`A` et `B`).

## Compilation et utilisation

Depuis la racine du dépôt, lancez le programme avec Go :

```sh
go run checker/checker.go 3 2 1
```

Les nombres peuvent être fournis séparément ou dans un même argument :

```sh
go run checker/checker.go "3 2 1"
```

Les instructions sont lues une par ligne depuis l'entrée standard. Par
exemple, avec PowerShell :

```powershell
"sa" | go run .\checker\checker.go 2 3
```

Pour vérifier la sortie d'un programme `push_swap`, transmettez ses
instructions au checker :

```sh
./push_swap 3 2 1 | go run checker/checker.go 3 2 1
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
