const fs = require("fs");
const path = require("path");
const csv = require("fast-csv");

async function main() {
  const result = { nodes: [], links: [] };

  fs.createReadStream(path.resolve(__dirname, "data", "dependents-bcgov.csv"))
    .pipe(csv.parse({ headers: true }))
    // pipe the parsed input into a csv formatter
    .pipe(csv.format({ headers: true }))
    // Using the transform function from the formatting stream
    .transform((row, next) => {
      if (row["Dependent Count"] !== "0") {
        result.nodes.push({ id: `bcgov/${row.Repository}`, group: 1 });

        const dependets = row.Dependents.split(";");

        dependets.forEach((target) =>
          result.links.push({ source: `bcgov/${row.Repository}`, target, value: 1 })
        );
      }

      next(null);
    })
    .on("finish", () => {
      fs.writeFile(
        path.resolve(
          __dirname,
          "data",
          "dependents-force-directed-graph-bcgov.json"
        ),
        JSON.stringify(result),
        "utf8",
        () => process.exit()
      );
    });
}

main();
