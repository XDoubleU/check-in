import { Migration } from '@mikro-orm/migrations';

export class Migration20230312202744 extends Migration {

  async up(): Promise<void> {
    this.addSql('create table "School" ("id" serial primary key, "name" varchar(255) not null);');
    this.addSql('alter table "School" add constraint "School_name_unique" unique ("name");');

    this.addSql('create table "User" ("id" uuid not null, "username" varchar(255) not null, "password_hash" varchar(255) not null, "roles" text[] not null default \'{user}\', constraint "User_pkey" primary key ("id"));');
    this.addSql('alter table "User" add constraint "User_username_unique" unique ("username");');

    this.addSql('create table "Location" ("id" uuid not null, "name" varchar(255) not null, "capacity" int not null, "user_id" uuid not null, constraint "Location_pkey" primary key ("id"), constraint Location_capacity_check check (capacity >= 0));');
    this.addSql('alter table "Location" add constraint "Location_user_id_unique" unique ("user_id");');

    this.addSql('create table "CheckIn" ("id" serial primary key, "location_id" uuid not null, "school_id" int not null default 1, "capacity" int not null, "created_at" timestamptz(0) not null);');

    this.addSql('alter table "Location" add constraint "Location_user_id_foreign" foreign key ("user_id") references "User" ("id") on update cascade;');

    this.addSql('alter table "CheckIn" add constraint "CheckIn_location_id_foreign" foreign key ("location_id") references "Location" ("id") on update cascade on delete cascade;');
    this.addSql('alter table "CheckIn" add constraint "CheckIn_school_id_foreign" foreign key ("school_id") references "School" ("id") on update cascade on delete set default;');
  }

  async down(): Promise<void> {
    this.addSql('alter table "CheckIn" drop constraint "CheckIn_school_id_foreign";');

    this.addSql('alter table "Location" drop constraint "Location_user_id_foreign";');

    this.addSql('alter table "CheckIn" drop constraint "CheckIn_location_id_foreign";');

    this.addSql('drop table if exists "School" cascade;');

    this.addSql('drop table if exists "User" cascade;');

    this.addSql('drop table if exists "Location" cascade;');

    this.addSql('drop table if exists "CheckIn" cascade;');
  }

}
