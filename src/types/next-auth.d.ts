/* eslint-disable @typescript-eslint/no-unused-vars */
import NextAuth from "next-auth"
import { JWT } from "next-auth/jwt"

declare module "next-auth" {
  interface Session {
    user: {
      id: string;
      username: string;
      isAdmin: boolean;
    };
  }

  interface User {
    id: string;
    username: string;
    isAdmin: boolean;
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    user: {
      id?: string;
      username?: string;
      isAdmin?: boolean
    };
  }
}
