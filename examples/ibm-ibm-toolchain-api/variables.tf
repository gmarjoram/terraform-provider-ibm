variable "ibmcloud_api_key" {
  description = "IBM Cloud API key"
  type        = string
}

// Resource arguments for toolchain_tool_pipeline
variable "toolchain_tool_pipeline_toolchain_id" {
  description = ""
  type        = string
  default     = "toolchain_id"
}
variable "toolchain_tool_pipeline_parameters_references" {
  description = "Decoded values used on provision in the broker that reference fields in the parameters."
  type        = map()
  default     = { "key": null }
}
