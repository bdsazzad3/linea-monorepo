/*
 * load simulation - OpenAPI 3.0
 * describe list of requests
 *
 * The version of the OpenAPI document: 1.0.0
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package net.consensys.zkevm.load.swagger;

import java.util.Objects;

import com.google.gson.annotations.SerializedName;

import java.io.IOException;

import com.google.gson.JsonElement;

import java.util.HashSet;

import net.consensys.zkevm.load.model.JSON;

/**
 * Scenario
 */

public class Scenario {
  public static final String SERIALIZED_NAME_TRANSACTION_TYPE = "scenarioType";
  @SerializedName(SERIALIZED_NAME_TRANSACTION_TYPE)
  protected String scenarioType;

  public Scenario() {
    this.scenarioType = this.getClass().getSimpleName();
  }

  public Scenario scenarioType(String scenarioType) {
    this.scenarioType = scenarioType;
    return this;
  }

   /**
   * Get scenarioType
   * @return scenarioType
  **/
  @javax.annotation.Nonnull
  public String getScenarioType() {
    return scenarioType;
  }

  public void setScenarioType(String scenarioType) {
    this.scenarioType = scenarioType;
  }



  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    Scenario scenario = (Scenario) o;
    return Objects.equals(this.scenarioType, scenario.scenarioType);
  }

  @Override
  public int hashCode() {
    return Objects.hash(scenarioType);
  }

  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class Scenario {\n");
    sb.append("    scenarioType: ").append(toIndentedString(scenarioType)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }


  public static HashSet<String> openapiFields;
  public static HashSet<String> openapiRequiredFields;

  static {
    // a set of all properties/fields (JSON key names)
    openapiFields = new HashSet<String>();
    openapiFields.add("scenarioType");

    // a set of required properties/fields (JSON key names)
    openapiRequiredFields = new HashSet<String>();
    openapiRequiredFields.add("scenarioType");
  }

 /**
  * Validates the JSON Element and throws an exception if issues found
  *
  * @param jsonElement JSON Element
  * @throws IOException if the JSON Element is invalid with respect to Scenario
  */
  public static void validateJsonElement(JsonElement jsonElement) throws IOException {
      if (jsonElement == null) {
        if (!Scenario.openapiRequiredFields.isEmpty()) { // has required fields but JSON element is null
          throw new IllegalArgumentException(String.format("The required field(s) %s in Scenario is not found in the empty JSON string", Scenario.openapiRequiredFields.toString()));
        }
      }

      String discriminatorValue = jsonElement.getAsJsonObject().get("scenarioType").getAsString();
      switch (discriminatorValue) {
        case "ContractCall":
          ContractCall.validateJsonElement(jsonElement);
          break;
        case "RoundRobinMoneyTransfer":
          RoundRobinMoneyTransfer.validateJsonElement(jsonElement);
          break;
        case "SelfTransactionWithPayload":
          SelfTransactionWithPayload.validateJsonElement(jsonElement);
          break;
        case "SelfTransactionWithRandomPayload":
          SelfTransactionWithRandomPayload.validateJsonElement(jsonElement);
          break;
        case "UnderPricedTransaction":
          UnderPricedTransaction.validateJsonElement(jsonElement);
          break;
        default:
          throw new IllegalArgumentException(String.format("The value of the `scenarioType` field `%s` does not match any key defined in the discriminator's mapping.", discriminatorValue));
      }
  }


 /**
  * Create an instance of Scenario given an JSON string
  *
  * @param jsonString JSON string
  * @return An instance of Scenario
  * @throws IOException if the JSON string is invalid with respect to Scenario
  */
  public static Scenario fromJson(String jsonString) throws IOException {
    return JSON.getGson().fromJson(jsonString, Scenario.class);
  }

 /**
  * Convert an instance of Scenario to an JSON string
  *
  * @return JSON string
  */
  public String toJson() {
    return JSON.getGson().toJson(this);
  }
}
