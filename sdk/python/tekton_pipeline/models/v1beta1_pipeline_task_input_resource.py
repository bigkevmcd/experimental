# Copyright 2020 The Tekton Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# coding: utf-8

"""
    Tekton

    Tekton Pipeline  # noqa: E501

    The version of the OpenAPI document: v0.17.2
    Generated by: https://openapi-generator.tech
"""


import pprint
import re  # noqa: F401

import six

from tekton_pipeline.configuration import Configuration


class V1beta1PipelineTaskInputResource(object):
    """NOTE: This class is auto generated by OpenAPI Generator.
    Ref: https://openapi-generator.tech

    Do not edit the class manually.
    """

    """
    Attributes:
      openapi_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    openapi_types = {
        '_from': 'list[str]',
        'name': 'str',
        'resource': 'str'
    }

    attribute_map = {
        '_from': 'from',
        'name': 'name',
        'resource': 'resource'
    }

    def __init__(self, _from=None, name=None, resource=None, local_vars_configuration=None):  # noqa: E501
        """V1beta1PipelineTaskInputResource - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration()
        self.local_vars_configuration = local_vars_configuration

        self.__from = None
        self._name = None
        self._resource = None
        self.discriminator = None

        if _from is not None:
            self._from = _from
        self.name = name
        self.resource = resource

    @property
    def _from(self):
        """Gets the _from of this V1beta1PipelineTaskInputResource.  # noqa: E501

        From is the list of PipelineTask names that the resource has to come from. (Implies an ordering in the execution graph.)  # noqa: E501

        :return: The _from of this V1beta1PipelineTaskInputResource.  # noqa: E501
        :rtype: list[str]
        """
        return self.__from

    @_from.setter
    def _from(self, _from):
        """Sets the _from of this V1beta1PipelineTaskInputResource.

        From is the list of PipelineTask names that the resource has to come from. (Implies an ordering in the execution graph.)  # noqa: E501

        :param _from: The _from of this V1beta1PipelineTaskInputResource.  # noqa: E501
        :type: list[str]
        """

        self.__from = _from

    @property
    def name(self):
        """Gets the name of this V1beta1PipelineTaskInputResource.  # noqa: E501

        Name is the name of the PipelineResource as declared by the Task.  # noqa: E501

        :return: The name of this V1beta1PipelineTaskInputResource.  # noqa: E501
        :rtype: str
        """
        return self._name

    @name.setter
    def name(self, name):
        """Sets the name of this V1beta1PipelineTaskInputResource.

        Name is the name of the PipelineResource as declared by the Task.  # noqa: E501

        :param name: The name of this V1beta1PipelineTaskInputResource.  # noqa: E501
        :type: str
        """
        if self.local_vars_configuration.client_side_validation and name is None:  # noqa: E501
            raise ValueError("Invalid value for `name`, must not be `None`")  # noqa: E501

        self._name = name

    @property
    def resource(self):
        """Gets the resource of this V1beta1PipelineTaskInputResource.  # noqa: E501

        Resource is the name of the DeclaredPipelineResource to use.  # noqa: E501

        :return: The resource of this V1beta1PipelineTaskInputResource.  # noqa: E501
        :rtype: str
        """
        return self._resource

    @resource.setter
    def resource(self, resource):
        """Sets the resource of this V1beta1PipelineTaskInputResource.

        Resource is the name of the DeclaredPipelineResource to use.  # noqa: E501

        :param resource: The resource of this V1beta1PipelineTaskInputResource.  # noqa: E501
        :type: str
        """
        if self.local_vars_configuration.client_side_validation and resource is None:  # noqa: E501
            raise ValueError("Invalid value for `resource`, must not be `None`")  # noqa: E501

        self._resource = resource

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.openapi_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, V1beta1PipelineTaskInputResource):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, V1beta1PipelineTaskInputResource):
            return True

        return self.to_dict() != other.to_dict()