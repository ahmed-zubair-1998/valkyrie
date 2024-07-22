resource "aws_lb" "frontend_server_alb" {
  name_prefix        = "fe-"
  load_balancer_type = "application"
  security_groups    = [aws_security_group.frontend_alb.id]
  subnets            = aws_subnet.public.*.id
  idle_timeout       = 60
  ip_address_type    = "dualstack"

  tags = { "Name" = "${var.default_tags.project}-frontend-server-alb" }
}

resource "aws_lb_target_group" "frontend_server_alb_targets" {
  name_prefix          = "fe-"
  port                 = 8080
  protocol             = "HTTP"
  vpc_id               = aws_vpc.main.id
  deregistration_delay = 30
  target_type          = "instance"

  health_check {
    enabled             = true
    path                = "/heartbeat"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 30
    interval            = 60
    protocol            = "HTTP"
  }

  tags = { "Name" = "${var.default_tags.project}-frontend-server-tg" }
}

resource "aws_lb_target_group_attachment" "frontend_server" {
  count = var.frontend_server_count

  target_group_arn = aws_lb_target_group.frontend_server_alb_targets.arn
  target_id        = aws_instance.frontend_server[count.index].id
  port             = 8080
}

resource "aws_lb_listener" "frontend_server_alb_http_80" {
  load_balancer_arn = aws_lb.frontend_server_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.frontend_server_alb_targets.arn
  }
}
