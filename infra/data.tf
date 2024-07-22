data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_instance" "dispatcher" {
  instance_id = aws_instance.dispatcher.id
}
