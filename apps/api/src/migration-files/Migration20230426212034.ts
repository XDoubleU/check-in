import { Migration } from '@mikro-orm/migrations';

export class Migration20230426212034 extends Migration {

  async up(): Promise<void> {
    this.addSql('alter table "Location" drop constraint Location_capacity_check;');
  }

  async down(): Promise<void> {
    this.addSql('alter table "Location" add constraint Location_capacity_check check(capacity >= 0);');
  }

}
