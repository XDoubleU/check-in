import {
  type PipeTransform,
  Injectable,
  type ArgumentMetadata,
  BadRequestException
} from "@nestjs/common"

@Injectable()
export class ParseDatePipe implements PipeTransform<string, Date> {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public transform(value: string, _metadata: ArgumentMetadata): Date {
    const val = new Date(value)
    if (isNaN(val.valueOf())) {
      throw new BadRequestException("Validation failed (Date is expected)")
    }
    return val
  }
}
