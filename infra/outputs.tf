output "client_alb_dns" {
  value       = aws_lb.frontend_server_alb.dns_name
  description = "DNS name of the AWS ALB for FE servers"
}

output "dispatcher_dns" {
  value       = aws_instance.dispatcher.public_dns
  description = "Public DNS of Dispatcher"
}
