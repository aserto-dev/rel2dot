# rel2dot

Create GraphViz dot file from Topaz relations JSON file.

## install rel2dot

```
go install github.com/aserto-dev/rel2dot@latest
```

## generate dot file
```
rel2dot -i ./relations.json | dot -Tpng > gdrive.png
```

## visualize dot file

![result](./gdrive.png)

## generate dot file with inverted directionality

```
rel2dot -f -i ./relations.json | dot -Tpng > gdrive-inv.png
```
## visualize dot file (inverted)

![result](./gdrive-inv.png)
