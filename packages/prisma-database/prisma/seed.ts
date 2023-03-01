import { PrismaClient } from "@prisma/client"

const prisma = new PrismaClient()

async function main() {
  const others = await prisma.school.upsert({
    where: { name: "Andere" },
    update: {},
    create: {
      name: "Andere"
    }
  })
  console.log({ others })
}
main()
  .then(async () => {
    await prisma.$disconnect()
  })
  .catch(async (e) => {
    console.error(e)
    await prisma.$disconnect()
    process.exit(1)
  })