# push-swap

Programme en Go qui trie une liste d'entiers en utilisant deux piles (`a` et `b`) et un jeu d'instructions limité. Il affiche la suite d'instructions, une par ligne, qui permet de trier la pile `a` par ordre croissant (le plus petit nombre au sommet).

## Prérequis

- Go (version indiquée dans [go.mod](go.mod))

## Utilisation

Depuis le dossier [push-swap/](push-swap/) :

```bash
go run . 2 1 3 6 5 8
```

Les nombres peuvent aussi être passés dans une seule chaîne :

```bash
go run . "2 1 3 6 5 8"
```

Exemple de sortie :

```
ra
pb
pb
pb
ra
ra
pa
pa
ra
pb
pb
ra
pa
pa
ra
pb
ra
pa
pa
```

Comportement :

| Cas | Résultat |
| --- | --- |
| Aucun argument | Aucune sortie |
| Pile déjà triée | Aucune sortie |
| Argument non entier ou doublon | `Error` sur la sortie d'erreur |

## Instructions disponibles

| Instruction | Effet |
| --- | --- |
| `sa` / `sb` / `ss` | Échange les deux premiers éléments de `a` / `b` / des deux piles |
| `pa` | Pousse le premier élément de `b` sur `a` |
| `pb` | Pousse le premier élément de `a` sur `b` |
| `ra` / `rb` / `rr` | Rotation vers le haut (le premier devient le dernier) |
| `rra` / `rrb` / `rrr` | Rotation vers le bas (le dernier devient le premier) |

## Algorithme

Le choix de la stratégie dépend du nombre d'éléments ([sort.go](push-swap/sort.go)) :

- **2 éléments** : un `sa`.
- **3 éléments** : tri par cas, 2 instructions maximum.
- **4 ou 5 éléments** : les plus petits éléments sont envoyés dans `b` par le chemin le plus court (`ra` ou `rra`), les 3 restants sont triés, puis tout est ramené avec `pa`.
- **6 éléments ou plus** : tri **radix** binaire. Les valeurs sont d'abord normalisées en rangs (0 à n-1), puis, bit par bit du poids faible au poids fort, les éléments dont le bit vaut 0 passent dans `b` (`pb`) et les autres tournent (`ra`), avant de tout ramener dans `a`.

Une étape d'optimisation finale (`optimize`) :

- supprime les instructions qui s'annulent (`pa`/`pb`, `ra`/`rra`, `rb`/`rrb`) ;
- fusionne les paires compatibles (`ra`+`rb` → `rr`, `rra`+`rrb` → `rrr`, `sa`+`sb` → `ss`).

## Structure du projet

```
push-swap/
├── go.mod
└── push-swap/
    ├── main.go       # point d'entrée
    ├── utils.go      # parsing des arguments, vérification du tri
    ├── ops.go        # piles et instructions
    ├── sort.go       # algorithmes de tri et optimisation
    └── sort_test.go  # tests
```

## Tests

```bash
cd push-swap
go test ./...
```
