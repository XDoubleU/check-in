import NextAuth, { NextAuthOptions } from "next-auth"
import { PrismaAdapter } from "@next-auth/prisma-adapter"
import { PrismaClient } from "@prisma/client"
import CredentialsProvider from "next-auth/providers/credentials"
import { compareSync } from "bcrypt"
import { User } from "next-auth/core/types"

// Instantiate Prisma Client
const prisma = new PrismaClient()

export const authOptions: NextAuthOptions = {
  session: {
    strategy: "jwt"
  },
  providers: [
    CredentialsProvider({
      name: "Credentials",
      credentials: {
        username: { label: "Username", type: "text", placeholder: "username" },
        password: { label: "Password", type: "password" }
      },
      async authorize(credentials) {
        if (credentials === undefined) {
          return null
        }

        try {
          const user = await prisma.user.findFirst({
            where: {
              username: credentials.username
            }
          })

          if (user === null) {
            return null
          }

          if (!compareSync(credentials.password, user.passwordHash)) {
            return null
          }
          
          const sessionUser : User = user
          if (!sessionUser.isAdmin) {
            const location = await prisma.location.findFirst({
              where: {
                userId: sessionUser.id
              }
            })

            if (location === null) {
              return null
            }
            
            sessionUser.locationId = location.id
          }

          return sessionUser
        }
        catch (err: unknown) {
          throw new Error("Authorize error: ", (err as Error))
        }
      }
    })
  ],
  pages: {
    signIn: "../../auth/signin"
  },
  callbacks: {
    async session({ session, token }) {
      session.user = {
        id: token.user.id ?? "",
        username: token.user.username ?? "",
        isAdmin: token.user.isAdmin ?? false,
        locationId: token.user.locationId
      }
      return Promise.resolve(session)
    },
    async jwt({ token, user }) {
      if (user) {
        token.user = user
      }
      return Promise.resolve(token)
    }
  },
  adapter: PrismaAdapter(prisma),
}

export default NextAuth(authOptions)